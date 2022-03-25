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

func Unmarshal(data []byte, v *Data, opts map[string]any) error {
	if err := jce.Unmarshal(data, v, true); err != nil {
		return err
	}
	switch v.Version {
	case 1, 2:
		p := make(Payload)
		if err := jce.Unmarshal(v.Payload, p); err != nil {
			return err
		}
		for key, b1 := range p {
			for _, b2 := range b1 {
				if opt, ok := opts[key]; ok {
					if err := jce.Unmarshal(b2, opt); err != nil {
						return err
					}
				}
			}
		}
	case 3:
		p := make(PayloadV3)
		if err := jce.Unmarshal(v.Payload, p); err != nil {
			return err
		}
		for key, b := range p {
			if opt, ok := opts[key]; ok {
				if err := jce.Unmarshal(b, opt); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
