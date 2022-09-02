package handlers

import (
	"fmt"
	"net/http"

	"github.com/getupio-undistro/zora/payloads"
	"github.com/getupio-undistro/zora/pkg/clientset/versioned"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ClusterListHandler(client versioned.Interface, logger logr.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.WithName("handlers.clusters").WithValues("method", r.Method, "path", r.URL.Path)
		clusterList, err := client.ZoraV1alpha1().Clusters("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Error(err, "failed to list Clusters")
			RespondWithDetailedError(w, http.StatusInternalServerError, "Error listing Clusters", err.Error())
			return
		}
		scanList, err := client.ZoraV1alpha1().ClusterScans("").List(r.Context(), metav1.ListOptions{})
		if err != nil {
			log.Error(err, "failed to list ClusterScans")
			RespondWithDetailedError(w, http.StatusInternalServerError, "Error listing ClusterScans", err.Error())
			return
		}

		clusters := payloads.NewClusterSlice(clusterList.Items, scanList.Items)
		log.Info(fmt.Sprintf("%d cluster(s) returned", len(clusters)))
		RespondWithJSON(w, http.StatusOK, clusters)
	}
}
