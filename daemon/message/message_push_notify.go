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
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/daemon/service"
	"github.com/elap5e/penguin/pkg/encoding/uni"
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
	data, push := uni.Data{}, PushNotifyRequest{}
	if err := uni.Unmarshal(reply.Payload[4:], &data, map[string]any{
		"req_PushNotify": &push,
	}); err != nil {
		return nil, err
	}
	items := []*service.DeleteMessage{}
	m.setSyncFlag(reply.Uin, 0)
	for force := true; force || m.getSyncFlag(reply.Uin) == 1; force = false {
		resp, err := m.GetMessage(reply.Uin)
		if err != nil {
			return nil, err
		}
		for _, uniPairMsgs := range resp.GetUinPairMsgs() {
			for _, msg := range uniPairMsgs.GetMsg() {
				head := msg.GetMsgHead()
				switch head.GetMsgType() {
				case 9, 10, 31, 79, 97, 120, 132, 133, 141, 166, 167:
					switch head.GetC2CCmd() {
					case 11, 175:
						if err := m.d.OnRecvMessage(head, msg.GetMsgBody()); err != nil {
							return nil, err
						}
						items = append(items, &service.DeleteMessage{
							FromUin:  int64(head.GetFromUin()),
							Time:     int64(head.GetMsgTime()),
							Seq:      int16(head.GetMsgSeq()),
							Cookie:   []byte{},
							Cmd:      int16(head.GetC2CCmd()),
							Type:     int64(head.GetMsgType()),
							AppID:    int64(head.GetFromAppid()),
							SendTime: 0,
							SSOSeq:   0,
							SSOIP:    0,
							ClientIP: 0,
						})
					}
				case 78, 81, 103, 107, 110, 111, 114, 118:
					_, _ = m.DeleteMessage(reply.Uin, &pb.MsgService_PbDeleteMsgReq_MsgItem{
						FromUin: head.GetFromUin(),
						ToUin:   head.GetToUin(),
						MsgType: head.GetMsgType(),
						MsgSeq:  head.GetMsgSeq(),
						MsgUid:  head.GetMsgUid(),
						Sig:     nil,
					})
				}
			}
		}
	}
	resp := service.OnlinePushMessageResponse{
		Uin:      reply.Uin,
		Items:    items,
		ServerIP: push.ServerIP,
		Token:    nil,
		Type:     0,
		Device:   nil,
	}
	return m.d.GetServiceManager().ResponseOnlinePushMessage(reply, &resp)
}
