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

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/sq325/k8sExporters/pkg/collector"
	"github.com/sq325/k8sExporters/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	versionInfo = "v1.3.0"
)

var (
	// /var/lib/kubelet/pods/uid/volumes/kubernetes.io~empty-dir
	port       *string = pflag.StringP("port", "p", "0", "listening port")
	prefixPath *string = pflag.String("prefixPath", "/var/lib/kubelet/pods", "the path of emptydir mount in node")
	version    *bool   = pflag.BoolP("version", "v", false, "Version info")
)

func main() {
	pflag.Parse()
	if *version {
		fmt.Println(versionInfo)
		return
	}

	// service account
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	if config == nil {
		log.Fatal("serviceAccount is nil")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	podfactor := resource.NewPodFactor(clientset)

	pesCollector, err := collector.NewEmptydirCollector(podfactor, *prefixPath)
	if err != nil {
		log.Fatal("NewEmptydirCollector Error:", err)
	}
	PromRegister(pesCollector)

	// http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>Emptydir Exporter</title></head>
				<body>
				<h1>emptydir usage</h1>
				<p>please click <a href="` + "metrics" + `">Metrics</a></p>
				</body>
				</html>`))
	})
	http.Handle("/metrics", promhttp.Handler())
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Address:", listener.Addr().String())
	_port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	log.Println("Listening port:", _port)
	log.Println("Metrics Url: http://<ip>:" + _port + "/metrics")
	log.Fatal(http.Serve(listener, nil))
}

func PromRegister(c prometheus.Collector) {
	prometheus.Unregister(promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}))
	prometheus.Unregister(promcollectors.NewGoCollector())
	prometheus.Register(c)
}
