#!/bin/sh
# Copyright 2022 Undistro Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit

KCONFIG_NAME=${KCONFIG_NAME:-"zora_view_kubeconfig.yaml"}
CLUSTER_ROLE_NAME=${CLUSTER_ROLE_NAME:-"zora-view"}
SVC_ACCOUNT_NS=${SVC_ACCOUNT_NS:-"zora-system"}
SVC_ACCOUNT_NAME=${SVC_ACCOUNT_NAME:-"zora-view"}
SVC_ACCOUNT_SECRET_NS=${SVC_ACCOUNT_SECRET_NS:-$SVC_ACCOUNT_NS}
SVC_ACCOUNT_SECRET_NAME=${SVC_ACCOUNT_SECRET_NAME:-"$SVC_ACCOUNT_NAME-token"}
METRICS_SERVER_VERSION=${METRICS_SERVER_VERSION:-"latest"}
METRICS_SERVER_DEPLOYMENT_NAME=${METRICS_SERVER_DEPLOYMENT_NAME:-"metrics-server"}
METRICS_SERVER_DEPLOYMENT=${METRICS_SERVER_DEPLOYMENT:-"https://github.com/kubernetes-sigs/metrics-server/releases/$METRICS_SERVER_VERSION/download/components.yaml"}

get_token_name() {
	echo $(kubectl -n $SVC_ACCOUNT_NS \
		get serviceaccount $SVC_ACCOUNT_NAME \
		-o jsonpath='{.secrets[0].name}'
	)
}
get_token_value() {
	echo $(kubectl -n $SVC_ACCOUNT_NS \
		get secret $TOKEN_NAME \
		-o jsonpath='{.data.token}' | base64 --decode
	)
}
get_current_context() {
	echo $(kubectl config current-context)
}
get_current_cluster() {
	echo $(kubectl config view \
		--raw -o go-template='
			{{range .contexts}}
				{{if eq .name "'$CURRENT_CONTEXT'"}}
					{{index .context "cluster"}}
				{{end}}
			{{end}}
		'
	)
}
get_cluster_ca() {
	echo $(kubectl config view \
		--raw -o go-template='
			{{range .clusters}}
				{{if eq .name "'$CURRENT_CLUSTER'"}}
					{{with index .cluster "certificate-authority-data"}}
						{{.}}
					{{end}}
				{{end}}
			{{end}}
		'
	)
}
get_cluster_server() {
	echo $(kubectl config view \
		--raw -o go-template='
			{{range .clusters}}
				{{if eq .name "'$CURRENT_CLUSTER'"}}
					{{ .cluster.server }}
				{{end}}
			{{end}}
		'
	)
}

create_svc_account() {
	kubectl -n $SVC_ACCOUNT_NS create serviceaccount $SVC_ACCOUNT_NAME
}

create_svc_account_secret() {
cat << EOF | kubectl create -f -
apiVersion: v1
kind: Secret
metadata:
  name: "$SVC_ACCOUNT_SECRET_NAME"
  namespace: "$SVC_ACCOUNT_SECRET_NS"
  annotations:
    kubernetes.io/service-account.name: "$SVC_ACCOUNT_NAME"
type: kubernetes.io/service-account-token
EOF
}

create_cluster_role() {
cat << EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: $CLUSTER_ROLE_NAME
rules:
  - apiGroups: [ "" ]
    resources:
      - configmaps
      - endpoints
      - limitranges
      - namespaces
      - nodes
      - persistentvolumes
      - persistentvolumeclaims
      - pods
      - replicationcontrollers
      - secrets
      - serviceaccounts
      - services
    verbs: [ "get", "list" ]
  - apiGroups: [ "apps" ]
    resources:
      - daemonsets
      - deployments
      - statefulsets
      - replicasets
    verbs: [ "get", "list" ]
  - apiGroups: [ "autoscaling" ]
    resources:
      - horizontalpodautoscalers
    verbs: [ "get", "list" ]
  - apiGroups: [ "networking.k8s.io" ]
    resources:
      - ingresses
      - networkpolicies
    verbs: [ "get", "list" ]
  - apiGroups: [ "policy" ]
    resources:
      - poddisruptionbudgets
      - podsecuritypolicies
    verbs: [ "get", "list" ]
  - apiGroups: [ "rbac.authorization.k8s.io" ]
    resources:
      - clusterroles
      - clusterrolebindings
      - roles
      - rolebindings
    verbs: [ "get", "list" ]
  - apiGroups: [ "metrics.k8s.io" ]
    resources:
      - pods
      - nodes
    verbs: [ "get", "list" ]
  - apiGroups: [ batch ]
    resources:
      - jobs
      - cronjobs
    verbs: [ "get", "list" ]
  - apiGroups: [ admissionregistration.k8s.io ]
    resources:
      - validatingwebhookconfigurations
      - mutatingwebhookconfigurations
    verbs: [ "get", "list" ]
EOF
}

create_cluster_role_binding() {
	kubectl create clusterrolebinding $SVC_ACCOUNT_NAME \
		--clusterrole=zora-view \
		--serviceaccount=$SVC_ACCOUNT_NS:$SVC_ACCOUNT_NAME
}

create_kubeconfig() {
cat << EOF > $KCONFIG_NAME
apiVersion: v1
kind: Config
current-context: $CURRENT_CONTEXT
contexts:
  - name: $CURRENT_CONTEXT
    context:
      cluster: $CURRENT_CONTEXT
      user: $SVC_ACCOUNT_NAME
      namespace: $SVC_ACCOUNT_NS
clusters:
  - name: $CURRENT_CONTEXT
    cluster:
      certificate-authority-data: $CLUSTER_CA
      server: $CLUSTER_SERVER
users:
  - name: $SVC_ACCOUNT_NAME
    user:
      token: $TOKEN_VALUE
EOF
}

setup_metrics_server() {
	if ! kubectl get pods -A 2> /dev/null | grep -q $METRICS_SERVER_DEPLOYMENT_NAME; then
		kubectl apply -f "$METRICS_SERVER_DEPLOYMENT"
	fi
}

setup_namespaces() {
	if ! kubectl get namespace $SVC_ACCOUNT_NS > /dev/null 2>&1; then
		kubectl create namespace $SVC_ACCOUNT_NS
	fi
}
setup_svc_account() {
	if ! kubectl -n $SVC_ACCOUNT_NS get serviceaccount $SVC_ACCOUNT_NAME > /dev/null 2>&1; then
		create_svc_account
	fi
}
setup_svc_account_secret() {
	if ! kubectl -n $SVC_ACCOUNT_SECRET_NS get secret $SVC_ACCOUNT_SECRET_NAME > /dev/null 2>&1; then
		create_svc_account_secret
	fi
}
setup_cluster_role() {
	if ! kubectl -n $SVC_ACCOUNT_NS get clusterrole $CLUSTER_ROLE_NAME > /dev/null 2>&1; then
		create_cluster_role
	fi
}
setup_cluster_role_binding() {
	if ! kubectl get -n $SVC_ACCOUNT_NS clusterrolebinding $SVC_ACCOUNT_NAME > /dev/null 2>&1; then
		create_cluster_role_binding
	fi
}


setup_namespaces
setup_svc_account

if kubectl version --short | awk '/Server/{if ($3 < "1.24.0") {exit 1}}'; then
  setup_svc_account_secret
  TOKEN_NAME=${TOKEN_NAME:-"$SVC_ACCOUNT_SECRET_NAME"}
else
  TOKEN_NAME=${TOKEN_NAME:-"$(get_token_name)"}
fi

TOKEN_VALUE=${TOKEN_VALUE:-"$(get_token_value)"}
CURRENT_CONTEXT=${CURRENT_CONTEXT:-"$(get_current_context)"}
CURRENT_CLUSTER=${CURRENT_CLUSTER:-"$(get_current_cluster)"}
CLUSTER_CA=${CLUSTER_CA:-"$(get_cluster_ca)"}
CLUSTER_SERVER=${CLUSTER_SERVER:-"$(get_cluster_server)"}

setup_metrics_server
setup_cluster_role
setup_cluster_role_binding
create_kubeconfig
