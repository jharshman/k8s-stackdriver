/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"errors"
	"strings"
	"sync"

	gce "cloud.google.com/go/compute/metadata"
)

var Config *GceConfig
var once sync.Once

// GceConfig aggregates all GCE related configuration parameters.
type GceConfig struct {
	Project         string
	Zone            string
	Cluster         string
	ClusterLocation string
	// This is actually instance name.
	// instance/name is not available from the GKE metadata server.
	// in that case it is replaced with instance/hostname.
	Instance   string
	InstanceId string
}

// GetGceConfig builds GceConfig based on the provided prefix and metadata server available on GCE.
func GetGceConfig() error {
	once.Do(func() {
		project, err := gce.ProjectID()

		cluster, err := gce.InstanceAttributeValue("cluster-name")
		cluster = strings.TrimSpace(cluster)

		// instance/name endpoint is not available on the GKE metadata server.
		// Try GCE instance/name endpoint. If error, try instance/hostname.
		// If instance/hostname, remove domain to replicate instance/name.
		node, err := gce.InstanceName()
		if err != nil {
			node, err = gce.Hostname()
			if err == nil {
				node = strings.Split(node, ".")[0]
			}
		}

		instanceId, _ := gce.InstanceID()
		zone, _ := gce.Zone()
		clusterLocation, _ := gce.InstanceAttributeValue("cluster-location")
		clusterLocation = strings.TrimSpace(clusterLocation)

		Config = &GceConfig{
			Project:         project,
			Zone:            zone,
			Cluster:         cluster,
			ClusterLocation: clusterLocation,
			Instance:        node,
			InstanceId:      instanceId,
		}
	})

	// ensure that Configuration is good
	if Config.Project == "" {
		return errors.New("project ID empty")
	}
	if Config.Cluster == "" {
		return errors.New("cluster empty")
	}
	if Config.ClusterLocation == "" {
		return errors.New("cluster location empty")
	}
	if Config.Instance == "" {
		return errors.New("instance/host name empty")
	}
	if Config.InstanceId == "" {
		return errors.New("instance ID empty")
	}
	if Config.Zone == "" {
		return errors.New("zone empty")
	}
	return nil
}
