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
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/channel/pb"
	pbmsg "github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (m *Manager) handlePushMessage(reply *rpc.Reply) (*rpc.Args, error) {
	push := pb.MsgPush_MsgOnlinePush{}
	if err := proto.Unmarshal(reply.Payload, &push); err != nil {
		return nil, err
	}
	for _, msg := range push.GetMsgs() {
		typ, sub := msg.GetHead().GetContentHead().GetMsgType(), msg.GetHead().GetContentHead().GetSubType()
		if typ == 3840 {
			_ = m.onChannelMessage(reply.Uin, msg)
		} else if typ == 3841 {
			_ = m.onChannelService(reply.Uin, msg)
		} else {
			dumpUnknown(typ, sub, msg)
		}
	}
	return nil, nil
}

func (m *Manager) onChannelMessage(uin int64, msg *pb.Common_Msg) error {
	return m.OnRecvChannelMessage(uin, msg)
}

func (m *Manager) onChannelService(uin int64, msg *pb.Common_Msg) error {
	typ, sub := msg.GetHead().GetContentHead().GetMsgType(), msg.GetHead().GetContentHead().GetSubType()
	var got *pbmsg.IMMsgBody_CommonElem
	for _, elem := range msg.GetBody().GetRichText().GetElems() {
		if elem := elem.GetCommonElem(); elem != nil {
			got = elem
			break
		}
	}
	if sub == 43 { // ChannelMute
		var notify pb.Channel_Service
		if err := proto.Unmarshal(got.GetPbElem(), &notify); err != nil {
			log.Error("failed to unmarshal channel event: %v", err)
			return nil
		}
		if v := notify.GetChannelMute(); v != nil {
			return m.onChannelServiceChannelMute(uin, v)
		}
	}
	dumpUnknown(typ, sub, msg)
	var notify pb.ChannelService_EventBody
	if err := proto.Unmarshal(got.GetPbElem(), &notify); err != nil {
		log.Error("failed to unmarshal channel event: %v", err)
		return nil
	}
	p, _ := json.Marshal(&notify)
	log.Warn("unknown msg type:%d sub_type:%d body:%s", typ, sub, p)
	if v := notify.GetMsgCreateGuild(); v != nil {
		return m.onChannelServiceChannelCreate(uin, v)
	} else if v := notify.GetMsgChangeGuildInfo(); v != nil {
		return m.onChannelServiceChannelUpdate(uin, v)
	} else if v := notify.GetMsgDestroyGuild(); v != nil {
		return m.onChannelServiceChannelDelete(uin, v)
	} else if v := notify.GetMsgCreateChan(); v != nil {
		return m.onChannelServiceRoomCreate(uin, v)
	} else if v := notify.GetMsgChangeChanInfo(); v != nil {
		return m.onChannelServiceRoomUpdate(uin, v)
	} else if v := notify.GetMsgDestroyChan(); v != nil {
		return m.onChannelServiceRoomDelete(uin, v)
	} else if v := notify.GetMsgCommGrayTips(); v != nil {
		return m.onChannelServiceGrayTips(uin, v)
	} else if v := notify.GetMsgUpdateMsg(); v != nil {
		return m.onChannelServiceMessageUpdate(uin, v)
	} else if v := notify.GetMsgJoinGuild(); v != nil {
		return m.onChannelServiceUserJoined(uin, v)
	} else if v := notify.GetMsgKickOffGuild(); v != nil {
		return m.onChannelServiceUserKicked(uin, v)
	} else if v := notify.GetMsgQuitGuild(); v != nil {
		return m.onChannelServiceUserLeaved(uin, v)
	}
	return nil
}

func (m *Manager) onChannelServiceMessageUpdate(uin int64, msg *pb.ChannelService_UpdateMsg) error {
	return nil
}

func (m *Manager) onChannelServiceChannelCreate(uin int64, msg *pb.ChannelService_CreateGuild) error {
	return nil
}

func (m *Manager) onChannelServiceChannelUpdate(uin int64, msg *pb.ChannelService_ChangeGuildInfo) error {
	return nil
}

func (m *Manager) onChannelServiceChannelDelete(uin int64, msg *pb.ChannelService_DestroyGuild) error {
	return nil
}

func (m *Manager) onChannelServiceChannelMute(uin int64, msg *pb.Channel_Service_ChannelMute) error {
	channel, _ := m.GetChannel(msg.GetChannelId())
	to, _ := m.GetUser(msg.GetChannelId(), msg.GetTinyId())
	toDisplay := to.Display
	if toDisplay == "" {
		toDisplay = to.Account.Username
	}
	if msg.GetTime() == 0 {
		log.Notifyf("[%d] channel:%d(%s) [%d(%s)]被解除禁言", uin, channel.ID, channel.Title, to.Account.ID, toDisplay)
	} else {
		log.Notifyf("[%d] channel:%d(%s) [%d(%s)]被禁言至[%s]", uin, channel.ID, channel.Title, to.Account.ID, toDisplay, time.Unix(msg.GetTime(), 0).UTC().Format(time.RFC3339))
	}
	return nil
}

func (m *Manager) onChannelServiceRoomCreate(uin int64, msg *pb.ChannelService_CreateChan) error {
	return nil
}

func (m *Manager) onChannelServiceRoomUpdate(uin int64, msg *pb.ChannelService_ChangeChanInfo) error {
	return nil
}

func (m *Manager) onChannelServiceRoomDelete(uin int64, msg *pb.ChannelService_DestroyChan) error {
	return nil
}

func (m *Manager) onChannelServiceGrayTips(uin int64, msg *pb.ChannelService_CommGrayTips) error {
	return nil
}

func (m *Manager) onChannelServiceUserJoined(uin int64, msg *pb.ChannelService_JoinGuild) error {
	return nil
}

func (m *Manager) onChannelServiceUserKicked(uin int64, msg *pb.ChannelService_KickOffGuild) error {
	return nil
}

func (m *Manager) onChannelServiceUserLeaved(uin int64, msg *pb.ChannelService_QuitGuild) error {
	return nil
}

func dumpUnknown(typ, sub uint64, msg *pb.Common_Msg) {
	p, _ := json.Marshal(msg)
	log.Warn("unknown msg type:%d sub_type:%d msg:%s", typ, sub, p)
}
