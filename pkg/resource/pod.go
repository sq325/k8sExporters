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

// PodFactor is only a single instance

package resource

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Service implement Resource interface
// Pod define a pod resource
type Pod struct {
	name      string
	namespace string
	kind      string // title style
	labels    map[string]string
}

// NewPod create a new pod
func NewPod(name, namespace string, labels map[string]string) *Pod {
	return &Pod{
		name:      name,
		namespace: namespace,
		kind:      "Pod",
		labels:    labels,
	}
}

func (p *Pod) Name() string {
	return p.name
}

func (p *Pod) Namespace() string {
	return p.namespace
}

func (p *Pod) Kind() string {
	return p.kind
}

func (p *Pod) Labels() map[string]string {
	return p.labels
}

func (p *Pod) Selector() map[string]string {
	return nil
}

type Pods []*Pod

// PodFactor implements Factor interface
// PodFactor parse output and produce Pods
type PodFactor struct {
	ClientSet *kubernetes.Clientset
}

func NewPodFactor(clientSet *kubernetes.Clientset) Factor {
	return &PodFactor{ClientSet: clientSet}
}

func (p *PodFactor) GetResources() (Resources, error) {
	podList, err := p.ClientSet.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	var pods Resources
	for _, npod := range podList.Items {
		pod := NewPod(npod.Name, npod.Namespace, npod.Labels)
		pods = append(pods, pod)
	}

	if len(pods) == 0 {
		log.Println("No pods found in cluster")
		return nil, nil
	}
	return pods, nil
}
