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

package service

import (
	"github.com/elap5e/penguin/pkg/encoding/jce"
)

type Message0x0210 struct {
	SubType int64  `jce:"0" json:"sub_type"`
	Payload []byte `jce:"10" json:"payload"`
}

func (m *Manager) decode0x0210Jce(uin int64, p []byte) (any, error) {
	msg := Message0x0210{}
	if err := jce.Unmarshal(p, &msg, true); err != nil {
		return nil, err
	}
	return m.decode0x0210(uin, msg.SubType, msg.Payload)
}

func (m *Manager) decode0x0210Pb(uin int64, p []byte) (any, error) {
	return nil, nil
}

func (m *Manager) decode0x0210(uin, typ int64, p []byte) (any, error) {
	return nil, nil
}
