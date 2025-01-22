/*
 * Copyright 2021 VMware, Inc.
 * All Rights Reserved.
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*   http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
package apiserver

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/api/models"

	"github.com/vmware/global-load-balancing-services-for-kubernetes/gslb/gslbutils"
	"github.com/vmware/global-load-balancing-services-for-kubernetes/gslb/nodes"
)

type GSGraphAPI struct{}

func (g GSGraphAPI) InitModel() {}

func (g GSGraphAPI) ApiOperationMap(prometheusEnavbled bool, reg *prometheus.Registry) []models.OperationMap {
	get := models.OperationMap{
		Route:   "/api/gsgraph",
		Method:  "GET",
		Handler: GSGraphHandler,
	}
	return []models.OperationMap{get}
}

func GSGraphHandler(w http.ResponseWriter, r *http.Request) {
	agl := nodes.SharedAviGSGraphLister()

	names, ok := r.URL.Query()["name"]
	if ok {
		tenant, exists := r.URL.Query()["tenant"]
		var tenName string
		if !exists {
			tenName = gslbutils.GetTenant() + "/" + names[0]
		} else {
			tenName = tenant[0] + "/" + names[0]
		}
		_, aviGS := agl.Get(tenName)
		if aviGS == nil {
			WriteToResponse(w, nil)
			return
		}
		aviGSGraph := aviGS.(*nodes.AviGSObjectGraph)
		WriteToResponse(w, aviGSGraph)
		return
	}

	keys := agl.GetAll()
	results := []interface{}{}

	for _, key := range keys {
		_, aviGS := agl.Get(key)
		aviGSGraph := aviGS.(*nodes.AviGSObjectGraph)
		results = append(results, aviGSGraph.GetCopy())
	}
	WriteToResponse(w, results)
}
