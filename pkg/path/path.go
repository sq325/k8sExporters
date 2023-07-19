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

package path

import (
	"os"
	"path/filepath"
)

// Represent file path or dir path
type Path struct {
	AbsPath  string
	IsDir    bool
	fileSize int64
}

func NewPath(name string) (*Path, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	abspath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	isDir := info.IsDir()
	return &Path{
		AbsPath:  abspath,
		IsDir:    isDir,
		fileSize: info.Size(),
	}, nil
}

// size of path
// if p is a file，return the size of file
// if p is dir, 遍历文件夹，统计所有文件和文件夹大小
// unit: byte
func (p *Path) Size() int64 {
	if p.IsDir {
		return p.dirUsage()
	}
	return p.fileSize
}

// unit: byte
func (p *Path) dirUsage() int64 {
	var size int64
	filepath.WalkDir(p.AbsPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 统计size
		info, err := d.Info()
		if err != nil {
			return err
		}
		size += info.Size()

		return nil
	})

	return size
}
