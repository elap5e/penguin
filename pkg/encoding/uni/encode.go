// Copyright 2022 Elapse and contributors
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

package uni

import (
	"github.com/elap5e/penguin/pkg/encoding/jce"
)

func Marshal(v *Data, opts map[string]any) ([]byte, error) {
	var err error
	switch v.Version {
	case 1, 2:
		p := make(Payload)
		for key, opt := range opts {
			b, err := jce.Marshal(opt)
			if err != nil {
				return nil, err
			}
			p[key][key] = b // TODO: fix
		}
		v.Payload, err = jce.Marshal(p)
		if err != nil {
			return nil, err
		}
	case 3:
		p := make(PayloadV3)
		for key, opt := range opts {
			b, err := jce.Marshal(opt)
			if err != nil {
				return nil, err
			}
			p[key] = b
		}
		v.Payload, err = jce.Marshal(p)
		if err != nil {
			return nil, err
		}
	}
	return jce.Marshal(v, true)
}
