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

func (m *Manager) DeleteMessage(uin int64, items ...*pb.MsgService_PbDeleteMsgReq_MsgItem) (*pb.MsgService_PbDeleteMsgResp, error) {
	if len(items) == 0 {
		return &pb.MsgService_PbDeleteMsgResp{}, nil
	}
	req := pb.MsgService_PbDeleteMsgReq{MsgItems: items}
	p, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.Call(service.MethodMessageDeleteMessage, &args, &reply); err != nil {
		return nil, err
	}
	resp := pb.MsgService_PbDeleteMsgResp{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
