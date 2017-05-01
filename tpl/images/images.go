// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package images

import (
	"errors"
	"image"
	"sync"

	// Importing image codecs for image.DecodeConfig
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/deps"
)

// New returns a new instance of the images-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		cache: map[string]image.Config{},
		deps:  deps,
	}
}

// Namespace provides template functions for the "images" namespace.
type Namespace struct {
	sync.RWMutex
	cache map[string]image.Config

	deps *deps.Deps
}

// Config returns the image.Config for the specified path relative to the
// working directory.
func (ns *Namespace) Config(path interface{}) (image.Config, error) {
	filename, err := cast.ToStringE(path)
	if err != nil {
		return image.Config{}, err
	}

	if filename == "" {
		return image.Config{}, errors.New("config needs a filename")
	}

	// Check cache for image config.
	ns.RLock()
	config, ok := ns.cache[filename]
	ns.RUnlock()

	if ok {
		return config, nil
	}

	f, err := ns.deps.Fs.WorkingDir.Open(filename)
	if err != nil {
		return image.Config{}, err
	}

	config, _, err = image.DecodeConfig(f)

	ns.Lock()
	ns.cache[filename] = config
	ns.Unlock()

	return config, err
}