// Copyright 2023 sunquan
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

	"github.com/spf13/pflag"
	"github.com/sq325/k8sExporters/pkg/nodevolume"
	"github.com/sq325/k8sExporters/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	versionInfo = "v1.0"
)

var (
	// /var/lib/kubelet/pods/uid/volumes/kubernetes.io~empty-dir
	prefixPath *string = pflag.String("prefixPath", "/var/lib/kubelet/pods", "the path of emptydir mount in node")
	version    *bool   = pflag.BoolP("version", "v", false, "Version info")
)

func main() {
	pflag.Parse()
	if *version {
		fmt.Println(versionInfo)
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
	podlist, err := podfactor.GetResources()
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range podlist {
		podemptydir, err := nodevolume.NewPodEmptyDir(pod.(*resource.Pod), *prefixPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(podemptydir.Pod.Name(), podemptydir.EmptyDir.SizeByte())
	}
}
