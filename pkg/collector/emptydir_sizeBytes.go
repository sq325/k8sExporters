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
	"github.com/sq325/k8sExporters/pkg/emptydir"
	"github.com/sq325/k8sExporters/pkg/resource"
)

// pod_emptydir_sizeBytes metric meta info
var (
	EmptydirSizeMetricName = "pod_emptydir_sizeBytes"
	Help                   = "pod_emptydir_sizeBytes displays the size of emptydir in pod"
	LabelKeys              = []string{"namespace", "pod", "hostIP", "emptydirName"}
)

var (
	emptydirSize_gaugeVec *prometheus.GaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: EmptydirSizeMetricName,
		Help: Help,
	}, LabelKeys)
)

// Emptydircollector implement prometheus.Collector interface
type EmptydirCollector struct {
	podEmptydirList []emptydir.IPodEmptydir
	gv              *prometheus.GaugeVec
}

// podList is a list of pod in this node
func NewEmptydirCollector(podList []*resource.Pod, factor emptydir.IPodEmptydir, prefixPath string) (*EmptydirCollector, error) {
	var podEmptydirList []emptydir.IPodEmptydir
	for _, pod := range podList {
		podEmptydir, err := emptydir.NewPodEmptydir(pod, prefixPath)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		podEmptydirList = append(podEmptydirList, podEmptydir)
	}

	return &EmptydirCollector{
		podEmptydirList: podEmptydirList,
		gv:              emptydirSize_gaugeVec,
	}, nil
}

func (e *EmptydirCollector) Describe(ch chan<- *prometheus.Desc) {
	e.gv.MetricVec.Describe(ch)
}

func (e *EmptydirCollector) Collect(ch chan<- prometheus.Metric) {
	for _, pod_emptydir := range e.podEmptydirList {
		e.gv.WithLabelValues(pod_emptydir.PodNamespace(), pod_emptydir.PodName(), pod_emptydir.PodHostIP()).Set(float64(pod_emptydir.EmptydirSizeBytes()))
	}
	e.gv.Collect(ch)
}
