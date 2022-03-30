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

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/channel/pb"
	"github.com/elap5e/penguin/pkg/encoding/pgn"
	"github.com/elap5e/penguin/pkg/log"
)

// 0x02000007ffffffff
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

func (d *Daemon) mustGetChannel(id int64, name string) *penguin.Chat {
	return &penguin.Chat{
		ID:    id,
		Type:  penguin.ChatTypeChannel,
		Title: name,
	}
}

func (d *Daemon) mustGetChannelRoom(cid, rid int64, extra *pb.Common_ExtInfo) *penguin.Chat {
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

func (d *Daemon) mustGetChannelRoomPrivate(cid, rid int64, name string) *penguin.Chat {
	return nil
}

func (d *Daemon) mustGetChannelUser(cid, uid int64, extra *pb.Common_ExtInfo) *penguin.User {
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
		// room
		channel := d.mustGetChannel(int64(rhead.GetGuildId()), string(extra.GetGuildName()))
		msg.Chat = d.mustGetChannelRoom(channel.ID, int64(rhead.GetChannelId()), extra)
		msg.From = d.mustGetChannelUser(channel.ID, int64(rhead.GetFromTinyid()), extra)
	} else {
		// room private
	}
	_ = pgn.NewMessageEncoder(recv.GetBody()).Encode(&msg)
	pi, _ := json.Marshal(recv)
	pm, _ := json.Marshal(msg)
	log.Debug("recv:%s msg:%s", pi, pm)
	log.Chat(id, &msg)
	return nil
}
