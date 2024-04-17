// Copyright 2023 Sun Quan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sq325/k8sExporters/pkg/resource"
	"github.com/sq325/k8sExporters/pkg/volumes"
)

// metrics meta info
var (
	sizeMetricName      = "pod_emptydir_size_bytes"
	sizeMetricHelp      = "pod_emptydir_size_bytes displays the size of emptydir in pod"
	sizeMetricLabelKeys = []string{"namespace", "pod", "hostIP", "emptydirName"}

	sizeLimitName      = "pod_emptydir_sizeLimit_bytes"
	sizeLimitHelp      = "pod_emptydir_sizeLimit_bytes displays the sizeLimit of emptydir in pod"
	sizeLimitLabelKeys = []string{"namespace", "pod", "hostIP", "emptydirName"}
)

var (
	sizeDesc *prometheus.Desc = prometheus.NewDesc(
		sizeMetricName,
		sizeLimitHelp,
		sizeMetricLabelKeys,
		nil,
	)
	sizeLimitDesc *prometheus.Desc = prometheus.NewDesc(
		sizeLimitName,
		sizeLimitHelp,
		sizeLimitLabelKeys,
		nil,
	)
)

// Emptydircollector implement prometheus.Collector interface
type EmptydirCollector struct {
	// /var/lib/kubelet/pods
	prefixPath string
	podfactor  *resource.PodFactor

	sizeDesc      *prometheus.Desc
	sizeLimitDesc *prometheus.Desc
}

// NewEmptydirCollector create a new emptydir collector for all pods in the cluster
func NewEmptydirCollector(factor *resource.PodFactor, prefixPath string) (*EmptydirCollector, error) {
	return &EmptydirCollector{
		prefixPath:    prefixPath,
		podfactor:     factor,
		sizeDesc:      sizeDesc,
		sizeLimitDesc: sizeLimitDesc,
	}, nil
}

func (e *EmptydirCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.sizeDesc
	ch <- e.sizeLimitDesc
}

func (e *EmptydirCollector) Collect(ch chan<- prometheus.Metric) {
	podemptydirList := e.getPodEmptydirList()
	for _, pod_emptydir := range podemptydirList {
		for vn, vs := range pod_emptydir.EmptydirListSizeBytes() {
			ch <- prometheus.MustNewConstMetric(e.sizeDesc, prometheus.GaugeValue, float64(vs), pod_emptydir.PodNamespace(), pod_emptydir.PodName(), pod_emptydir.PodHostIP(), vn)
		}
		for vn, vsl := range pod_emptydir.EmptydirListSizeLimitBytes() {
			ch <- prometheus.MustNewConstMetric(e.sizeLimitDesc, prometheus.GaugeValue, float64(vsl), pod_emptydir.PodNamespace(), pod_emptydir.PodName(), pod_emptydir.PodHostIP(), vn)
		}
	}
}

func (e *EmptydirCollector) getPodEmptydirList() []volumes.IPodEmptydirs {
	podlist, err := e.podfactor.GetPods() // all pods in the cluster
	if err != nil {
		log.Fatal(err)
	}

	var podEmptydirList []volumes.IPodEmptydirs
	for _, pod := range podlist {

		podEmptydirs, err := volumes.NewPodEmptydirs(pod, e.prefixPath)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
		podEmptydirList = append(podEmptydirList, podEmptydirs)
	}
	return podEmptydirList
}
