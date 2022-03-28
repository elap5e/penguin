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
	"encoding/hex"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/service/pb"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type OnlinePushMessageResponse struct {
	Uin      int64            `jce:"0" json:"uin"`
	Items    []*DeleteMessage `jce:"1" json:"items"`
	ServerIP uint32           `jce:"2" json:"server_ip"`
	Token    []byte           `jce:"3" json:"token"`
	Type     uint32           `jce:"4" json:"type"`
	Device   *Device          `jce:"5" json:"device,omitempty"`
}

type Device struct {
	NetworkType  uint8  `jce:"0" json:"network_type"`
	DeviceType   string `jce:"1" json:"device_type"`
	OSVersion    string `jce:"2" json:"os_version"`
	VendorName   string `jce:"3" json:"vendor_name"`
	VendorOSName string `jce:"4" json:"vendor_os_name"`
	IOSIDFA      string `jce:"5" json:"ios_idfa"`
}

type DeleteMessage struct {
	FromUin  int64  `jce:"0" json:"from_uin,omitempty"`
	Time     int64  `jce:"1" json:"time,omitempty"`
	Seq      int16  `jce:"2" json:"seq,omitempty"`
	Cookie   []byte `jce:"3" json:"cookie,omitempty"`
	Cmd      int16  `jce:"4" json:"method,omitempty"`
	Type     int64  `jce:"5" json:"type,omitempty"`
	AppID    int64  `jce:"6" json:"app_id,omitempty"`
	SendTime int64  `jce:"7" json:"send_time,omitempty"`
	SSOSeq   int32  `jce:"8" json:"sso_seq,omitempty"`
	SSOIP    int32  `jce:"9" json:"sso_ip,omitempty"`
	ClientIP int32  `jce:"10" json:"client_ip,omitempty"`
}

func (m *Manager) handleOnlinePushMessage(reply *rpc.Reply) (*rpc.Args, error) {
	log.Debug("handleOnlinePushMessage: \n%s", hex.Dump(reply.Payload))
	log.Debug("handleOnlinePushMessage: \n%s", hex.EncodeToString(reply.Payload))
	push := pb.OnlinePush_PbPushMsg{}
	if err := proto.Unmarshal(reply.Payload, &push); err != nil {
		return nil, err
	}
	head, body := push.GetMsg().GetMsgHead(), push.GetMsg().GetMsgBody()
	if err := m.d.OnRecvMessage(head, body); err != nil {
		return nil, err
	}
	return m.responseOnlinePushMessage(reply, &OnlinePushMessageResponse{
		Uin: reply.Uin,
		Items: []*DeleteMessage{{
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
		Token:    nil,
		Type:     0,
		Device:   nil,
	})
}

func (m *Manager) responseOnlinePushMessage(reply *rpc.Reply, resp *OnlinePushMessageResponse) (*rpc.Args, error) {
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
