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
	"encoding/json"

	"github.com/elap5e/penguin/daemon/message/dto"
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/daemon/service"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (m *Manager) handlePushNotifyRequest(reply *rpc.Reply) (*rpc.Args, error) {
	data, push := uni.Data{}, dto.PushNotifyRequest{}
	if err := uni.Unmarshal(reply.Payload[4:], &data, map[string]any{
		"req_PushNotify": &push,
	}); err != nil {
		return nil, err
	}
	items := []*dto.MessageDelete{}
	for {
		resp, err := m.GetMessage(reply.Uin)
		if err != nil {
			return nil, err
		}
		for _, uniPairMsgs := range resp.GetUinPairMsgs() {
			for _, msg := range uniPairMsgs.GetMsg() {
				head := msg.GetMsgHead()
				switch int32(head.GetMsgType()) {
				// case 9, 10, 31, 79, 97, 120, 132, 133, 141, 166, 167:
				case 9, 10, 31, 79, 97, 120, 132, 133, 166, 167: // 0x84(132)
					switch head.GetC2CCmd() {
					case 11, 175:
						if err := m.OnRecvMessage(reply.Uin, head, msg.GetMsgBody()); err != nil {
							return nil, err
						}
						items = append(items, &dto.MessageDelete{
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
					default:
						p, _ := json.Marshal(msg)
						log.Warn("unhandled private message cmd:%d msg:%s", head.GetC2CCmd(), p)
					}
				case 78, 81, 103, 107, 110, 111, 114, 118:
					_, _ = m.DeleteMessage(reply.Uin, &pb.MsgService_PbDeleteMsgReq_MsgItem{
						FromUin: head.GetFromUin(),
						ToUin:   head.GetToUin(),
						MsgType: head.GetMsgType(),
						MsgSeq:  head.GetMsgSeq(),
						MsgUid:  head.GetMsgUid(),
						Sig:     []byte{},
					})
				case 208:
					if err := m.OnRecvMessage(reply.Uin, head, msg.GetMsgBody()); err != nil {
						return nil, err
					}
				case 528: // 0x210(528)
					_ = m.GetServiceManager().Decode0x210(reply.Uin, &dto.Message{
						Type:         528,
						MessageBytes: msg.GetMsgBody().GetMsgContent(),
					}, true)
				case 732: // 0x2dc(732)
					_ = m.GetServiceManager().Decode0x2dc(reply.Uin, &dto.Message{
						Type:         732,
						MessageBytes: msg.GetMsgBody().GetMsgContent(),
					})
				case -1010, -1009, -1008, -1007, -1006:
					// service message
					fallthrough
				case -1023, -1022, -1021, -1020:
					// chat service message
					fallthrough
				case 33, 34:
					// chat user joined
					fallthrough
				case 35, 36, 37, 45, 46, 84, 85, 86, 87:
					// chat service message
					fallthrough
				case 187, 188, 189, 190, 191:
					// service message
					fallthrough
				default:
					p, _ := json.Marshal(msg)
					log.Warn("unhandled type cmd:%d msg:%s", head.GetMsgType(), p)
				}
			}
		}
		if flag, _ := m.GetFlag(reply.Uin); flag != 1 {
			_, _ = m.SetFlag(reply.Uin, 0) // clear flag
			break
		}
	}
	resp := service.OnlinePushResponse{
		Uin:      reply.Uin,
		Items:    items,
		ServerIP: push.ServerIP,
		Token:    []byte{},
		Type:     0,
		Device:   nil,
	}
	return m.GetServiceManager().OnlinePushResponse(reply, &resp)
}
