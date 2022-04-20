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
	"github.com/elap5e/penguin/daemon/message/dto"
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/encoding/jce"
)

type Message0x210 struct {
	SubType int32  `jce:"0" json:"sub_type"`
	Payload []byte `jce:"10" json:"payload"`
}

func (m *Manager) Decode0x210(uin int64, msg *dto.Message, isPb ...bool) error {
	dumpUnknown(msg.Type, msg)
	if len(isPb) > 0 && isPb[0] {
		return m.decode0x210Pb(uin, msg.MessageBytes)
	}
	return m.decode0x210Jce(uin, msg.MessageBytes)
}

func (m *Manager) decode0x210Jce(uin int64, p []byte) error {
	msg := Message0x210{}
	if err := jce.Unmarshal(p, &msg, true); err != nil {
		return err
	}
	return m.decode0x210(uin, msg.SubType, msg.Payload)
}

func (m *Manager) decode0x210Pb(uin int64, p []byte) error {
	msg := pb.MsgCommon_MsgType0X210{}
	if err := jce.Unmarshal(p, &msg, true); err != nil {
		return err
	}
	return m.decode0x210(uin, int32(msg.GetSubMsgType()), msg.GetMsgContent())
}

func (m *Manager) decode0x210(uin int64, typ int32, p []byte) error {
	return nil
}
