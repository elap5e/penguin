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
	"fmt"
	"math/rand"

	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
	"google.golang.org/protobuf/proto"
)

func (m *Manager) SendMessage(uin int64, req *pb.MsgService_PbSendMsgReq) (*pb.MsgService_PbSendMsgResp, error) {
	if req.GetMsgSeq() == 0 {
		return nil, fmt.Errorf("invalid msg seq")
	}
	if req.GetMsgRand() == 0 {
		req.MsgRand = rand.Uint32()
	}
	if len(req.GetSyncCookie()) == 0 {
		req.SyncCookie, _ = m.GetCookie(uin)
	}
	p, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.c.Call(service.MethodMessageSendMessage, &args, &reply); err != nil {
		return nil, err
	}
	resp := pb.MsgService_PbSendMsgResp{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	// resp.Result
	//     0: success
	//     1: ???
	//    16: elements (notFriend)
	//   120: elements (groupMute)
	//   241: ???
	// -3902: marketFace (vip/svip)
	// -4902: marketFace magic (vip/svip)
	//  5002: poke (vip/svip)
	return &resp, nil
}
