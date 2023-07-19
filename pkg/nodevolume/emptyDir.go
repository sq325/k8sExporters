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

package nodevolume

import (
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

type PodEmptyDir struct {
	Pod      *resource.Pod
	EmptyDir *EmptyDir
}

func NewPodEmptyDir(pod *resource.Pod, prefixPath string) (*PodEmptyDir, error) {
	uid := pod.UID()
	emptydir, err := NewEmptyDir(prefixPath, uid)
	if err != nil {
		return nil, err
	}
	return &PodEmptyDir{
		Pod:      pod,
		EmptyDir: emptydir,
	}, nil
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

// emptyDir 代表某一个node上所有的emptydir
// 由node上的pod的emptydir组成
// {prefixPath}/{uid}/{DefaultKubeletVolumesDirName}/{emptyDirPluginName}
type EmptyDir struct {
	VolumePath *path.Path
}

func NewEmptyDir(prefixPath, uid string) (*EmptyDir, error) {
	path, err := path.NewPath(filepath.Join(prefixPath, DefaultKubeletVolumesDirName, emptyDirPluginName))
	if err != nil {
		return nil, err
	}
	return &EmptyDir{
		VolumePath: path,
	}, nil
}

// /var/lib/kubelet/pods/uid/{DefaultKubeletVolumesDirName}/{emptyDirPluginName}
func (e *EmptyDir) Path() *path.Path {
	return e.VolumePath
}

// global var
func (e *EmptyDir) PluginName() string {
	return emptyDirPluginName
}

func (e *EmptyDir) SizeByte() int64 {
	return e.Path().Size()
}
