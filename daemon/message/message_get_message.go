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
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

func (m *Manager) GetMessage(uin int64) (*pb.MsgService_PbGetMsgResp, error) {
	flag, cookie := m.getSyncFlag(uin), m.getSyncCookie(uin)
	req := pb.MsgService_PbGetMsgReq{
		SyncFlag:           flag,
		SyncCookie:         cookie,
		RambleFlag:         0,
		LatestRambleNumber: 20,
		OtherRambleNumber:  3,
		OnlineSyncFlag:     1, // fix
		ContextFlag:        1,
		WhisperSessionId:   0,
		MsgReqType:         0, // fix
		PubaccountCookie:   nil,
		MsgCtrlBuf:         nil,
		ServerBuf:          nil,
	}
	p, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.c.Call(service.MethodMessageGetMessage, &args, &reply); err != nil {
		return nil, err
	}
	resp := pb.MsgService_PbGetMsgResp{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	m.setSyncFlag(uin, resp.GetSyncFlag())
	m.setSyncCookie(uin, resp.GetSyncCookie())
	return &resp, nil
}
