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

package message

import (
	"encoding/hex"
	"encoding/json"

	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type PushNotifyRequest struct {
	Uin         int64  `jce:"0" json:"uin,omitempty"`
	Type        uint8  `jce:"1" json:"type,omitempty"`
	Service     string `jce:"2" json:"service,omitempty"`
	Cmd         string `jce:"3" json:"cmd,omitempty"`
	Cookie      []byte `jce:"4" json:"cookie,omitempty"`
	MessageType uint16 `jce:"5" json:"message_type,omitempty"`
	UserActive  uint32 `jce:"6" json:"user_active,omitempty"`
	GeneralFlag uint32 `jce:"7" json:"general_flag,omitempty"`
	BindedUin   int64  `jce:"8" json:"binded_uin,omitempty"`

	Message       *Message `jce:"9" json:"message,omitempty"`
	ControlBuffer string   `jce:"10" json:"control_buffer,omitempty"`
	ServerBuffer  []byte   `jce:"11" json:"server_buffer,omitempty"`
	PingFlag      uint64   `jce:"12" json:"ping_flag,omitempty"`
	ServerIP      uint32   `jce:"13" json:"server_ip,omitempty"`
}

func (m *Manager) handlePushNotifyRequest(reply *rpc.Reply) (*rpc.Args, error) {
	log.Debug("handlePushNotifyRequest: \n%s", hex.Dump(reply.Payload))
	data, req := uni.Data{}, PushNotifyRequest{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"req_PushNotify": &req,
	}); err != nil {
		return nil, err
	}
	p, _ := json.MarshalIndent(req, "", "  ")
	log.Debug("handlePushNotifyRequest:\n%s", p)
	return nil, nil
}
