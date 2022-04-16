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

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/channel/pb"
	pbmsg "github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/encoding/pgn"
	"github.com/elap5e/penguin/pkg/log"
)

// mask 0x02000007ffffffff
func (d *Daemon) prefetchChannelAccount(id int64) error {
	account, ok := d.accm.GetAccount(id)
	if !ok {
		account = &penguin.Account{
			ID:   id,
			Type: penguin.AccountTypeChannel,
		}
		_, _ = d.accm.SetAccount(account.ID, account)
	}
	return nil
}

func (d *Daemon) getOrLoadChannel(id int64, name string) *penguin.Chat {
	return &penguin.Chat{
		ID:    id,
		Type:  penguin.ChatTypeChannel,
		Title: name,
	}
}

func (d *Daemon) getOrLoadChannelRoom(cid, rid int64, extra *pb.Common_ExtInfo) *penguin.Chat {
	return &penguin.Chat{
		ID:   rid,
		Type: penguin.ChatTypeRoomText,
		Chat: &penguin.Chat{
			ID:    cid,
			Type:  penguin.ChatTypeChannel,
			Title: string(extra.GetGuildName()),
		},
		Title: string(extra.GetChannelName()),
	}
}

func (d *Daemon) getOrLoadChannelRoomPrivate(cid, rid int64, name string) *penguin.Chat {
	return nil
}

func (d *Daemon) getOrLoadChannelUser(cid, uid int64, extra *pb.Common_ExtInfo) *penguin.User {
	return &penguin.User{
		Account: &penguin.Account{
			ID:       uid,
			Type:     penguin.AccountTypeChannel,
			Username: string(extra.GetMemberName()),
			Photo:    string(extra.GetFromAvatar()),
		},
		Display: string(extra.GetFromNick()),
	}
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
		cid, rid, fid := int64(rhead.GetGuildId()), int64(rhead.GetChannelId()), int64(rhead.GetFromUin())
		channel := d.getOrLoadChannel(cid, string(extra.GetGuildName()))
		msg.Chat = d.getOrLoadChannelRoom(channel.ID, rid, extra)
		msg.From = d.getOrLoadChannelUser(channel.ID, fid, extra)
	} else {
		// room private
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
	// identify message type
	if msg.Chat.Type == penguin.ChatTypeChannel {
		// channel
		return fmt.Errorf("a channel message can not be sent to a channel")
	} else if msg.Chat.Type == penguin.ChatTypeRoomText ||
		msg.Chat.Type == penguin.ChatTypeRoomVoice {
		// room any
		req.Msg.Head.RoutingHead = &pb.Common_RoutingHead{
			GuildId:   uint64(msg.Chat.Channel.ID),
			ChannelId: uint64(msg.Chat.ID),
			FromUin:   uint64(id),
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
	go d.pushMessage(msg)
	preq, _ := json.Marshal(req)
	prsp, _ := json.Marshal(resp)
	pmsg, _ := json.Marshal(msg)
	log.Debug("id:%d req:%s resp:%s msg:%s", id, preq, prsp, pmsg)
	log.Chat(id, msg)
	return nil
}
