{{/*
Expand the name of the chart.
*/}}
{{- define "zora.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "zora.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "zora.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "zora.labels" -}}
helm.sh/chart: {{ include "zora.chart" . }}
{{ include "zora.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Operator labels
*/}}
{{- define "zora.operatorLabels" -}}
{{ include "zora.labels" . }}
app.kubernetes.io/component: operator
{{- end }}

{{/*
Server labels
*/}}
{{- define "zora.serverLabels" -}}
{{ include "zora.labels" . }}
app.kubernetes.io/component: server
{{- end }}

{{/*
UI labels
*/}}
{{- define "zora.uiLabels" -}}
{{ include "zora.labels" . }}
app.kubernetes.io/component: ui
{{- end }}

{{/*
NGINX labels
*/}}
{{- define "zora.nginxLabels" -}}
{{ include "zora.labels" . }}
app.kubernetes.io/component: nginx
{{- end }}

{{/*
Selector labels
*/}}
{{- define "zora.selectorLabels" -}}
app.kubernetes.io/name: {{ include "zora.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Operator selector labels
*/}}
{{- define "zora.operatorSelectorLabels" -}}
{{ include "zora.selectorLabels" . }}
app.kubernetes.io/component: operator
{{- end }}

{{/*
Server selector labels
*/}}
{{- define "zora.serverSelectorLabels" -}}
{{ include "zora.selectorLabels" . }}
app.kubernetes.io/component: server
{{- end }}

{{/*
UI selector labels
*/}}
{{- define "zora.uiSelectorLabels" -}}
{{ include "zora.selectorLabels" . }}
app.kubernetes.io/component: ui
{{- end }}

{{/*
NGINX selector labels
*/}}
{{- define "zora.nginxSelectorLabels" -}}
{{ include "zora.selectorLabels" . }}
app.kubernetes.io/component: nginx
{{- end }}

{{/*
Create the name of the service account to use in Operator
*/}}
{{- define "zora.operatorServiceAccountName" -}}
{{- if .Values.operator.rbac.serviceAccount.create }}
{{- default (printf "%s-%s" (include "zora.fullname" .) "operator") .Values.operator.rbac.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.operator.rbac.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use in Server
*/}}
{{- define "zora.serverServiceAccountName" -}}
{{- if .Values.server.rbac.serviceAccount.create }}
{{- default (printf "%s-%s" (include "zora.fullname" .) "server") .Values.server.rbac.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.server.rbac.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use in UI
*/}}
{{- define "zora.uiServiceAccountName" -}}
{{- if .Values.ui.serviceAccount.create }}
{{- default (printf "%s-%s" (include "zora.fullname" .) "ui") .Values.ui.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.ui.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use in NGINX
*/}}
{{- define "zora.nginxServiceAccountName" -}}
{{- if .Values.nginx.serviceAccount.create }}
{{- default (printf "%s-%s" (include "zora.fullname" .) "nginx") .Values.nginx.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.nginx.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "imagePullSecret" }}
{{- with .Values.imageCredentials }}
{{- printf "{\"auths\":{\"%s\":{\"auth\":\"%s\"}}}" .registry (printf "%s:%s" .username .password | b64enc) | b64enc }}
{{- end }}
{{- end }}
