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

// PodFactor is only a single instance

package resource

import (
	"context"
	"log"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// Service implement Resource interface
// Pod define a pod resource
type Pod struct {
	pod           *coreV1.Pod
	namespaceName types.NamespacedName
	uid           types.UID
	kind          string // title style
	labels        map[string]string
}

// NewPod create a new pod
func NewPod(pod coreV1.Pod) *Pod {
	namespaceName := types.NamespacedName{
		Namespace: pod.Namespace,
		Name:      pod.Name,
	}

	return &Pod{
		pod:           &pod,
		namespaceName: namespaceName,
		uid:           pod.UID,
		kind:          "Pod",
		labels:        pod.Labels,
	}
}

func (p *Pod) Name() string {
	return p.namespaceName.Name
}

func (p *Pod) UID() string {
	return string(p.uid)
}

func (p *Pod) Namespace() string {
	return p.namespaceName.Namespace
}

func (p *Pod) HostIP() string {
	return p.pod.Status.HostIP
}

func (p *Pod) Volumes() []*Volume {
	var volumes []*Volume
	for _, volume := range p.pod.Spec.Volumes {
		volume := volume // notice
		if v := NewVolume(&volume); v.Type() != "" {
			volumes = append(volumes, v)
		}
	}
	return volumes
}

func (p *Pod) String() string {
	return p.namespaceName.String()
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

func NewPodFactor(clientSet *kubernetes.Clientset) *PodFactor {
	return &PodFactor{ClientSet: clientSet}
}

func (p *PodFactor) GetPods() (Pods, error) {
	podList, err := p.ClientSet.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	var pods Pods
	for _, npod := range podList.Items {
		pod := NewPod(npod)
		pods = append(pods, pod)
	}

	if len(pods) == 0 {
		log.Println("No pods found in cluster")
		return nil, nil
	}
	return pods, nil
}
