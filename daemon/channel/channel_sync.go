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

package channel

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/channel/pb"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

func (m *Manager) SyncFirstView(uin int64, seq uint32) (*pb.SyncLogic_FirstViewRsp, error) {
	return m.syncFirstView(uin, &pb.SyncLogic_FirstViewReq{
		LastMsgTime:       0,
		UdcFlag:           0,
		Seq:               seq,
		DirectMessageFlag: 1,
		NoNeedMsg:         1,
	})
}

func (m *Manager) syncFirstView(uin int64, req *pb.SyncLogic_FirstViewReq) (*pb.SyncLogic_FirstViewRsp, error) {
	p, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err := m.c.Call(service.MethodChannelSyncFirstView, &args, &reply); err != nil {
		return nil, err
	}
	resp := pb.SyncLogic_FirstViewRsp{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	p, _ = json.MarshalIndent(&resp, "", "  ")
	log.Debug("syncFirstView:\n%s", string(p))
	return &resp, nil
}

func (m *Manager) handlePushFirstView(reply *rpc.Reply) (*rpc.Args, error) {
	resp := pb.SyncLogic_FirstViewMsg{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	p, _ := json.MarshalIndent(&resp, "", "  ")
	log.Debug("handlePushFirstView:\n%s", string(p))
	if nodes := resp.GetDirectMessageGuildNodes(); nodes != nil {
		for _, node := range nodes {
			channel := penguin.Chat{
				ID:   int64(node.GetGuildId()),
				Type: penguin.ChatTypeChannel,
			}
			_, _ = m.SetChannel(channel.ID, &channel)
			p, _ := json.Marshal(channel)
			log.Debug("channel:%d:%s", channel.ID, p)
			for _, node := range node.GetChannelNodes() {
				room := penguin.Chat{
					ID:   int64(node.GetChannelId()),
					Type: penguin.ChatTypeRoomPrivate,
				}
				_, _ = m.SetRoom(channel.ID, room.ID, &room)
				p, _ := json.Marshal(room)
				log.Debug("channel:%d:room:%d:%s", channel.ID, room.ID, p)
			}
		}
	} else if nodes := resp.GetGuildNodes(); nodes != nil {
		for _, node := range nodes {
			channel := penguin.Chat{
				ID:    int64(node.GetGuildId()),
				Type:  penguin.ChatTypeChannel,
				Title: string(node.GetGuildName()),
			}
			_, _ = m.SetChannel(channel.ID, &channel)
			p, _ := json.Marshal(channel)
			log.Debug("channel:%d:%s", channel.ID, p)
			for _, node := range node.GetChannelNodes() {
				room := penguin.Chat{
					ID:    int64(node.GetChannelId()),
					Type:  penguin.ChatTypeRoomText,
					Title: string(node.GetChannelName()),
				}
				_, _ = m.SetRoom(channel.ID, room.ID, &room)
				p, _ := json.Marshal(room)
				log.Debug("channel:%d:room:%d:%s", channel.ID, room.ID, p)
			}
		}
	}
	return nil, nil
}
