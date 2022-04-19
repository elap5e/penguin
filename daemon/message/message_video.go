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

func (m *Manager) ChatDownloadVideo(uin int64) (*pb.PttShortVideo_RspBody, error) {
	req := &pb.PttShortVideo_ReqBody{
		Cmd: 400,
		Seq: 0,
	}
	req.Msg_PttShortVideoDownload_Req = &pb.PttShortVideo_PttShortVideoDownloadReq{
		Fromuin:                   0, // fix
		Touin:                     uint64(uin),
		ChatType:                  0,   // fix
		ClientType:                0,   // fix
		Fileid:                    "",  // fix
		GroupCode:                 0,   // fix
		FileMd5:                   nil, // fix
		AgentType:                 0,   // fix
		BusinessType:              0,   // fix
		FlagSupportLargeSize:      1,
		FlagClientQuicProtoEnable: 1,
		FileType:                  0, // fix
		DownType:                  0, // fix
		SceneType:                 0, // fix
		NeedInnerAddr:             0,
		ReqTransferType:           0, // fix
	}
	req.MsgExtensionReq = []*pb.PttShortVideo_ExtensionReq{{
		SubBusiType: 0, // fix
	}}
	return m.chatDownloadVideo(uin, req)
}

func (m *Manager) chatDownloadVideo(uin int64, req *pb.PttShortVideo_ReqBody) (*pb.PttShortVideo_RspBody, error) {
	p, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.Call(service.MethodMessageChatDownloadVideo, &args, &reply); err != nil {
		return nil, err
	}
	resp := pb.PttShortVideo_RspBody{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
