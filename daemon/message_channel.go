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

package daemon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/channel/pb"
	pbmsg "github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/encoding/pgn"
	"github.com/elap5e/penguin/pkg/log"
)

func (d *Daemon) prefetchChannelAccount(id int64) error {
	account, ok := d.accm.GetChannelAccount(id)
	if !ok {
		_, _ = d.accm.SetChannelAccount(account.ID, &penguin.Account{
			ID:   id,
			Type: penguin.AccountTypeChannel,
		})
	}
	return nil
}

func (d *Daemon) getOrLoadChannel(channelID int64, title string, private ...bool) *penguin.Chat {
	channel, ok := d.chnm.GetChannel(channelID)
	if !ok {
		typ := penguin.ChatTypeChannel
		if len(private) > 0 && private[0] {
			typ = penguin.ChatTypeChannelPrivate
		}
		_, _ = d.chnm.SetChannel(channelID, &penguin.Chat{
			ID:    channelID,
			Type:  typ,
			Title: title,
		})
		channel, _ = d.chnm.GetChannel(channelID)
	}
	return channel
}

func (d *Daemon) getOrLoadChannelRoom(channelID, roomID int64, ctrl *pb.Common_MsgCtrlHead, extra *pb.Common_ExtInfo) *penguin.Chat {
	room, ok := d.chnm.GetRoom(channelID, roomID)
	if !ok {
		typ, ctyp := penguin.ChatTypeRoomPrivate, ctrl.GetChannelType()
		channel := d.getOrLoadChannel(channelID, string(extra.GetGuildName()), ctyp == 0)
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
		_, _ = d.chnm.SetRoom(channelID, roomID, &penguin.Chat{
			ID:      roomID,
			Type:    typ,
			Title:   string(extra.GetChannelName()),
			Channel: channel,
		})
		room, _ = d.chnm.GetRoom(channelID, roomID)
	}
	if room.Channel.Type == penguin.ChatTypeChannelPrivate {
		// fix src channel
		srcChannelID, srcChannelTitle := int64(extra.GetDirectMessageMember()[0].GetSourceGuildId()), string(extra.GetDirectMessageMember()[0].GetSourceGuildName())
		room.Channel.Channel = d.getOrLoadChannel(srcChannelID, srcChannelTitle)
	}
	return room
}

func (d *Daemon) getOrLoadChannelUser(channelID, userID int64, extra *pb.Common_ExtInfo) *penguin.User {
	user, ok := d.chnm.GetUser(channelID, userID)
	if !ok {
		account, ok := d.accm.GetChannelAccount(userID)
		if !ok {
			_, _ = d.accm.SetChannelAccount(account.ID, &penguin.Account{
				ID:       userID,
				Type:     penguin.AccountTypeChannel,
				Username: string(extra.GetFromNick()),
				Photo:    string(extra.GetFromAvatar()),
			})
		}
		_, _ = d.chnm.SetUser(channelID, account.ID, &penguin.User{
			Account: account,
			Display: string(extra.GetMemberName()),
		})
		user, _ = d.chnm.GetUser(channelID, userID)
	}
	return user
}

func (d *Daemon) OnRecvChannelMessage(id int64, recv *pb.Common_Msg) error {
	chead, rhead := recv.GetHead().GetContentHead(), recv.GetHead().GetRoutingHead()
	extra := recv.GetExtInfo()
	msg := penguin.Message{
		MessageID: int64(chead.GetMsgSeq()),
		Time:      int64(chead.GetMsgTime()),
	}
	// pre-fetch accounts
	_ = d.prefetchChannelAccount(int64(rhead.GetFromTinyid()))
	if rhead.GetDirectMessageFlag() == 0 {
		// room any
		channelID, roomID, fromID := int64(rhead.GetGuildId()), int64(rhead.GetChannelId()), int64(rhead.GetFromTinyid())
		msg.Chat = d.getOrLoadChannelRoom(channelID, roomID, recv.GetCtrlHead(), extra)
		msg.From = d.getOrLoadChannelUser(channelID, fromID, extra)
	} else {
		// room private
		channelID, roomID, fromID := int64(rhead.GetGuildId()), int64(rhead.GetChannelId()), int64(rhead.GetFromTinyid())
		msg.Chat = d.getOrLoadChannelRoom(channelID, roomID, recv.GetCtrlHead(), extra)
		msg.From = d.getOrLoadChannelUser(channelID, fromID, extra)
	}
	if err := pgn.NewMessageEncoder(recv.GetBody()).Encode(&msg); err != nil {
		return err
	}
	return d.onRecvChannelMessage(id, recv, &msg)
}

func (d *Daemon) SendChannelMessage(id int64, msg *penguin.Message) error {
	var req pb.Oidb0Xf62_ReqBody
	req.Msg = &pb.Common_Msg{}
	req.Msg.Head = &pb.Common_MsgHead{}
	random := rand.New(rand.NewSource(time.Now().Unix()))
	// identify message type
	if msg.Chat.Type == penguin.ChatTypeChannel {
		// channel
		return fmt.Errorf("a channel message can not be sent to a channel")
	} else if msg.Chat.Type == penguin.ChatTypeRoomText ||
		msg.Chat.Type == penguin.ChatTypeRoomLive {
		// room any
		req.Msg.Head.RoutingHead = &pb.Common_RoutingHead{
			GuildId:   uint64(msg.Chat.Channel.ID),
			ChannelId: uint64(msg.Chat.ID),
			FromUin:   uint64(id),
		}
		req.Msg.Head.ContentHead = &pb.Common_ContentHead{
			MsgType: 3840,
			Random:  random.Uint64(),
		}
	} else if msg.Chat.Type == penguin.ChatTypeRoomPrivate {
		// room private
	} else {
		return fmt.Errorf("unknown chat type: %s", msg.Chat.Type)
	}
	// encode message
	req.Msg.CtrlHead = &pb.Common_MsgCtrlHead{}
	req.Msg.Body = &pbmsg.IMMsgBody_MsgBody{}
	if err := pgn.NewMessageDecoder(req.Msg.Body).Decode(msg); err != nil {
		return err
	}
	resp, err := d.chnm.SendMessage(id, &req)
	if err != nil {
		return err
	}
	return d.onSendChannelMessage(id, &req, resp, msg)
}

func (d *Daemon) onRecvChannelMessage(id int64, recv *pb.Common_Msg, msg *penguin.Message) error {
	go d.pushMessage(msg)
	go d.fetchBlobs(msg)
	pr, _ := json.Marshal(recv)
	pm, _ := json.Marshal(msg)
	log.Debug("id:%d recv:%s msg:%s", id, pr, pm)
	log.Chat(id, msg)
	return nil
}

func (d *Daemon) onSendChannelMessage(id int64, req *pb.Oidb0Xf62_ReqBody, resp *pb.Oidb0Xf62_RspBody, msg *penguin.Message) error {
	if msg.From.Account.Type == penguin.AccountTypeDefault {
		account, ok := d.accm.GetChannelFromDefault(msg.From.Account.ID)
		if ok {
			_ = msg.Chat.Channel.ID
			_ = account.ID
			msg.From = d.getOrLoadChannelUser(msg.Chat.Channel.ID, account.ID, nil)
		}
	}
	go d.pushMessage(msg)
	preq, _ := json.Marshal(req)
	prsp, _ := json.Marshal(resp)
	pmsg, _ := json.Marshal(msg)
	log.Debug("id:%d req:%s resp:%s msg:%s", id, preq, prsp, pmsg)
	log.Chat(id, msg)
	return nil
}
