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
	if err := m.Call(service.MethodChannelSyncFirstView, &args, &reply); err != nil {
		return nil, err
	}
	resp := pb.SyncLogic_FirstViewRsp{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	m.GetAccountManager().SetDefaultChannelPair(uin, int64(resp.GetSelfTinyid()))
	return &resp, nil
}

func (m *Manager) handlePushFirstView(reply *rpc.Reply) (*rpc.Args, error) {
	resp := pb.SyncLogic_FirstViewMsg{}
	if err := proto.Unmarshal(reply.Payload, &resp); err != nil {
		return nil, err
	}
	if nodes := resp.GetDirectMessageGuildNodes(); nodes != nil {
		for _, node := range nodes {
			if err := m.onPushFirstViewChannel(reply.Uin, penguin.ChatTypeChannelPrivate, node); err != nil {
				log.Println(err)
			}
		}
	} else if nodes := resp.GetGuildNodes(); nodes != nil {
		for _, node := range nodes {
			if err := m.onPushFirstViewChannel(reply.Uin, penguin.ChatTypeChannel, node); err != nil {
				log.Println(err)
			}
		}
	}
	return nil, nil
}

func (m *Manager) onPushFirstViewChannel(uin int64, typ penguin.ChatType, node *pb.SyncLogic_GuildNode) error {
	channelID := int64(node.GetGuildId())
	_, _ = m.SetChannel(channelID, &penguin.Chat{
		ID:    channelID,
		Type:  typ,
		Title: string(node.GetGuildName()),
	})
	channel, _ := m.GetChannel(channelID)
	p, _ := json.Marshal(channel)
	log.Debug("channel:%d:%s", channelID, p)
	for _, node := range node.GetChannelNodes() {
		typ, ctyp := penguin.ChatTypeRoomPrivate, node.GetChannelType()
		if ctyp == 0 {
		} else if ctyp == 1 {
			typ = penguin.ChatTypeRoomText
		} else if ctyp == 2 {
			typ = penguin.ChatTypeRoomVoice
		} else if ctyp == 4 {
			typ = penguin.ChatTypeRoomGroup
		} else if ctyp == 5 {
			typ = penguin.ChatTypeRoomLive
		} else if ctyp == 6 {
			typ = penguin.ChatTypeRoomApp
		} else if ctyp == 7 {
			typ = penguin.ChatTypeRoomForum
		} else {
			log.Warn("unknown channel type:%d", ctyp)
		}
		roomID := int64(node.GetChannelId())
		_, _ = m.SetRoom(channelID, roomID, &penguin.Chat{
			ID:      roomID,
			Type:    typ,
			Title:   string(node.GetChannelName()),
			Channel: channel,
		})
		room, _ := m.GetRoom(channelID, roomID)
		p, _ := json.Marshal(room)
		log.Debug("channel:%d:room:%d:%s", channelID, roomID, p)
	}
	if typ == penguin.ChatTypeChannel {
		var tinyID int64
		finish, offset, cookie := false, int64(0), []byte{}
		for !finish {
			resp, err := m.GetChannelUsers(uin, channelID, offset, 100, cookie)
			if err != nil {
				log.Println(err)
				break
			}
			finish, offset, cookie = resp.GetFinish() == 1, resp.GetOffset(), resp.GetCookie()
			tinyID = resp.GetOwner().GetTinyId()
		}
		_, err := m.GetChannelRoles(uin, channelID)
		if err != nil {
			log.Println(err)
		}
		if _, err := m.GetChannelUserRoles(uin, channelID, tinyID); err != nil {
			log.Println(err)
		}
	}
	return nil
}
