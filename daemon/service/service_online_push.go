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
	"encoding/json"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/message/dto"
	"github.com/elap5e/penguin/daemon/service/pb"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type OnlinePushRequest struct {
	Uin             int64                 `jce:",0" json:"uin,omitempty"`
	Time            int64                 `jce:",1" json:"time,omitempty"`
	Messages        []*dto.Message        `jce:",2" json:"messages,omitempty"`
	ServerIP        uint32                `jce:",3" json:"server_ip,omitempty"`
	SyncCookie      []byte                `jce:",4" json:"sync_cookie,omitempty"`
	UinPairMessages []*dto.UinPairMessage `jce:",5" json:"uin_pair_messages,omitempty"`
	Previews        map[string][]byte     `jce:",6" json:"previews,omitempty"`
	UserActive      int32                 `jce:",7" json:"user_active,omitempty"`
	GeneralFlag     int32                 `jce:",12" json:"general_flag,omitempty"`
}

type OnlinePushResponse struct {
	Uin      int64                `jce:"0" json:"uin"`
	Items    []*dto.MessageDelete `jce:"1" json:"items"`
	ServerIP uint32               `jce:"2" json:"server_ip"`
	Token    []byte               `jce:"3" json:"token"`
	Type     uint32               `jce:"4" json:"type"`
	Device   *Device              `jce:"5" json:"device,omitempty"`
}

type Device struct {
	NetworkType  uint8  `jce:"0" json:"network_type"`
	DeviceType   string `jce:"1" json:"device_type"`
	OSVersion    string `jce:"2" json:"os_version"`
	VendorName   string `jce:"3" json:"vendor_name"`
	VendorOSName string `jce:"4" json:"vendor_os_name"`
	IOSIDFA      string `jce:"5" json:"ios_idfa"`
}

func (m *Manager) handleOnlinePushMessage(reply *rpc.Reply) (*rpc.Args, error) {
	push := pb.OnlinePush_PbPushMsg{}
	if err := proto.Unmarshal(reply.Payload, &push); err != nil {
		return nil, err
	}
	head, body := push.GetMsg().GetMsgHead(), push.GetMsg().GetMsgBody()
	if err := m.d.OnRecvMessage(reply.Uin, head, body); err != nil {
		return nil, err
	}
	return m.OnlinePushResponse(reply, &OnlinePushResponse{
		Uin: reply.Uin,
		Items: []*dto.MessageDelete{{
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
		}},
		ServerIP: uint32(push.GetSvrip()),
		Token:    []byte{},
		Type:     0,
		Device:   nil,
	})
}

func (m *Manager) OnlinePushResponse(reply *rpc.Reply, resp *OnlinePushResponse) (*rpc.Args, error) {
	p, err := uni.Marshal(&uni.Data{
		Version:     3,
		RequestID:   reply.Seq,
		ServantName: "OnlinePush",
		FuncName:    "SvcRespPushMsg",
	}, map[string]any{
		"resp": resp,
	})
	if err != nil {
		return nil, err
	}
	return &rpc.Args{
		Version:       rpc.VersionSimple,
		Uin:           reply.Uin,
		Seq:           reply.Seq,
		ServiceMethod: service.MethodServiceOnlinePushResponse,
		Payload:       p,
	}, nil
}

func (m *Manager) handleOnlinePushRequest(reply *rpc.Reply) (*rpc.Args, error) {
	data, push := uni.Data{}, OnlinePushRequest{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]interface{}{
		"req": &push,
	}); err != nil {
		return nil, err
	}
	items := []*dto.MessageDelete{}
	for _, msg := range push.Messages {
		switch msg.Type {
		case 0x0210:
			// _, _ = m.decode0x0210Jce(reply.Uin, msg.MessageBytes)
			fallthrough
		case 0x02dc:
			// _, _ = m.decode0x02dc(reply.Uin, msg.MessageBytes)
			fallthrough
		default:
			p, _ := json.Marshal(msg)
			log.Warn("msg:%s", p)
		}
		items = append(items, &dto.MessageDelete{
			FromUin: msg.FromUin,
			Time:    msg.Time,
			Seq:     int16(msg.Seq),
			Cookie:  msg.MessageCookie,
		})
	}
	resp := OnlinePushResponse{
		Uin:      reply.Uin,
		Items:    items,
		ServerIP: push.ServerIP,
		Token:    []byte{},
		Type:     0,
		Device:   nil,
	}
	return m.OnlinePushResponse(reply, &resp)
}

func (m *Manager) handleOnlinePushTicketExpired(reply *rpc.Reply) (*rpc.Args, error) {
	if _, err := m.d.GetAuthManager().SignInChangeToken(reply.Uin); err != nil {
		return nil, err
	}
	if _, err := m.RegisterAppRegister(reply.Uin); err != nil {
		return nil, err
	}
	return &rpc.Args{
		Version:       rpc.VersionSimple,
		Uin:           reply.Uin,
		Seq:           reply.Seq,
		ServiceMethod: service.MethodServiceOnlinePushTicketExpired,
	}, nil
}
