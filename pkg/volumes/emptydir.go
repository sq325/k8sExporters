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

package volumes

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/sq325/k8sExporters/pkg/path"
	"github.com/sq325/k8sExporters/pkg/resource"
)

var (
	DefaultKubeletVolumesDirName = "volumes"
	emptyDirPluginName           = EscapeQualifiedName("kubernetes.io/empty-dir")
	// prefixPath                   = "/var/lib/kubelet/pods"
)

func EscapeQualifiedName(in string) string {
	return strings.Replace(in, "/", "~", -1)
}

type IPodEmptydir interface {
	PodName() string
	PodNamespace() string
	PodUID() string
	PodHostIP() string
	EmptydirListSizeBytes() map[string]int64
	EmptydirListSizeLimitBytes() map[string]int64
}

// PodEmptyDir implement PodEmptydirFactor interface
type PodEmptydir struct {
	Pod          *resource.Pod
	EmptydirList []*EmptyDir
}

func NewPodEmptydir(pod *resource.Pod, prefixPath string) (*PodEmptydir, error) {
	uid := pod.UID()

	var emptydirList []*EmptyDir
	for _, v := range pod.Volumes {
		emptydir, err := NewEmptyDir(prefixPath, uid, v)
		if err != nil {
			log.Println(err)
			continue
		}
		emptydirList = append(emptydirList, emptydir)
	}
	return &PodEmptydir{
		Pod:          pod,
		EmptydirList: emptydirList,
	}, nil
}

func (p *PodEmptydir) PodName() string {
	return p.Pod.Name()
}

func (p *PodEmptydir) PodNamespace() string {
	return p.Pod.Namespace()
}

func (p *PodEmptydir) PodUID() string {
	return p.Pod.UID()
}

func (p *PodEmptydir) PodHostIP() string {
	return p.Pod.HostIP()
}

func (p *PodEmptydir) EmptydirListSizeBytes() map[string]int64 {
	if len(p.EmptydirList) == 0 {
		return nil
	}
	m := make(map[string]int64)
	for _, e := range p.EmptydirList {
		m[e.volume.Name()] = e.SizeBytes()
	}
	return m
}

func (p *PodEmptydir) EmptydirListSizeLimitBytes() map[string]int64 {
	if len(p.EmptydirList) == 0 {
		return nil
	}
	m := make(map[string]int64)
	for _, e := range p.EmptydirList {
		m[e.volume.Name()] = e.volume.SizeLimit()
	}
	return m
}

// {prefixPath}/{uid}/{DefaultKubeletVolumesDirName}/{emptyDirPluginName}/{emptydirName}
type EmptyDir struct {
	path   *path.Path
	volume *resource.Volume
}

// NewEmptyDir代表此uid的pod volumes中的 名称为{emptydirName} 的emptydir
func NewEmptyDir(prefixPath, uid string, volume *resource.Volume) (*EmptyDir, error) {
	if volume.Type() != "EmptyDir" {
		return nil, errors.New(fmt.Sprint("pod", uid, "volume type is not EmptyDir"))
	}

	emptydirName := volume.Name()
	path, err := path.NewPath(filepath.Join(prefixPath, uid, DefaultKubeletVolumesDirName, emptyDirPluginName, emptydirName))
	if err != nil {
		return nil, err
	}
	return &EmptyDir{
		path:   path,
		volume: volume,
	}, nil
}

// /var/lib/kubelet/pods/uid/{DefaultKubeletVolumesDirName}/{emptyDirPluginName}
func (e *EmptyDir) Path() *path.Path {
	return e.path
}

// global var
func (e *EmptyDir) PluginName() string {
	return emptyDirPluginName
}

func (e *EmptyDir) SizeBytes() int64 {
	return e.Path().Size()
}

type NodeEmptydir struct {
}

// type EmptyDirList []*EmptyDir

// func NewEmptyDirList(podfactor resource.Factor, prefixPath string) (EmptyDirList, error) {
// 	// podList
// 	var errs error
// 	podlist, err := podfactor.GetResources()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// emptydirList
// 	emptyDirList := make(EmptyDirList, 0, len(podlist))

// 	// append emptydir to emptyDirList
// 	for _, pod := range podlist {
// 		uid := pod.(*(resource.Pod)).UID()
// 		emptydir, err := NewEmptyDir(prefixPath, uid)
// 		if err != nil {
// 			if errs == nil {
// 				errs = err
// 			} else {
// 				errs = errors.Join(errs, err)
// 			}
// 			continue
// 		}
// 		emptyDirList = append(emptyDirList, emptydir)
// 	}

// 	return emptyDirList, errs
// }
