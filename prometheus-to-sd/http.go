package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/k8s-stackdriver/prometheus-to-sd/config"
	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v3 "google.golang.org/api/monitoring/v3"
	"net/http"
)

func setupRouter() http.Server {

	probeStatus := &struct{
		StackDriver bool
		GCE *config.GceConfig
	}{
		StackDriver: stackdriverOK,
		GCE: config.Config,
	}

	router := chi.NewRouter()
	router.Get("/ready", func(writer http.ResponseWriter, request *http.Request) {
		err := json.NewEncoder(writer).Encode(probeStatus)
		if err != nil {
			glog.Warningf("failed to write probe status to endpoint: %v", err)
		}
		writer.WriteHeader(200)
		if probeStatus.StackDriver != true {
			writer.WriteHeader(500)
		}
	})
	router.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", *listenAddress, *metricsPort),
		Handler: router,
	}
	return srv
}

// ensureFunctional tests the basic functionality required for
// the StackDriver service. This will try listing, creating, and finally
// deleting a Metric Descriptor. This is done because it is possible for a
// monitoring v3 Service be created and NOT have sufficient permissions / scopes
// for proper function of prometheus-to-sd.
func ensureFunctional(stackDriverService *v3.Service) bool {
	// attempt to list all metric descriptors in project matching custom.google.api filter
	listReq := stackDriverService.Projects.MetricDescriptors.List(fmt.Sprintf("project/%s", config.Config.Project))
	listresp, err := listReq.Do()
	if err != nil {
		return false
	}
	if listresp.HTTPStatusCode != 200 {
		return false
	}

	// attempt to create a metric descriptor

	// attempt to delete the metric descriptor

	return true
}