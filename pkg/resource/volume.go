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

package resource

import coreV1 "k8s.io/api/core/v1"

type Volume struct {
	volume *coreV1.Volume
}

func NewVolume(volume *coreV1.Volume) *Volume {
	return &Volume{
		volume: volume,
	}
}

func (v *Volume) SizeLimit() int64 {
	var sizeLimit int64
	if v.Type() == "EmptyDir" {
		if v.volume.VolumeSource.EmptyDir.SizeLimit != nil {
			sizeLimit = v.volume.VolumeSource.EmptyDir.SizeLimit.Value()
		}
	}
	return sizeLimit
}

func (v *Volume) Name() string {
	return v.volume.Name
}

func (v *Volume) Type() string {
	var _type string
	switch {
	case v.volume.VolumeSource.EmptyDir != nil:
		_type = "EmptyDir"
	case v.volume.VolumeSource.HostPath != nil:
		_type = "HostPath"
	}
	return _type
}
