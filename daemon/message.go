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
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/encoding/pgn"
	"github.com/elap5e/penguin/pkg/log"
)

func (d *Daemon) OnRecvMessage(uin int64, head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) error {
	msg := penguin.Message{
		MessageID: int64(head.GetMsgSeq()),
		Time:      int64(head.GetMsgTime()),
	}
	msg.From = &penguin.User{
		Account: &penguin.Account{
			ID:   int64(head.GetFromUin()),
			Type: penguin.AccountTypeDefault,
		},
	}
	if v := head.GetDiscussInfo(); v != nil {
		// discuss
	} else if v := head.GetDiscussInfo(); v != nil {
		// discuss private
	} else if v := head.GetGroupInfo(); v != nil {
		// group
		from, ok := d.chtm.GetUser(int64(v.GetGroupCode()), int64(head.GetFromUin()))
		if !ok {
			account, ok := d.accm.GetAccount(int64(head.GetFromUin()))
			if !ok {
				account = &penguin.Account{
					ID:   int64(head.GetFromUin()),
					Type: penguin.AccountTypeDefault,
				}
				_, _ = d.accm.SetAccount(account.ID, account)
			}
			from = &penguin.User{
				Account: account,
			}
			_, _ = d.chtm.SetUser(int64(v.GetGroupCode()), account.ID, from)
		}
		msg.From = from
		chat, ok := d.chtm.GetChat(int64(v.GetGroupCode()))
		if !ok {
			chat = &penguin.Chat{
				ID:    int64(v.GetGroupCode()),
				Type:  penguin.ChatTypeGroup,
				Title: string(v.GetGroupName()),
			}
			_, _ = d.chtm.SetChat(int64(v.GetGroupCode()), chat)
		}
		msg.Chat = chat
	} else if v := head.GetC2CTmpMsgHead(); v != nil {
		// group private
		msg.Chat = &penguin.Chat{
			ID:   int64(v.GetGroupCode()),
			Type: penguin.ChatTypeGroupPrivate,
		}
		msg.Chat.User = &penguin.User{
			Account: &penguin.Account{
				ID:   int64(head.GetFromUin()),
				Type: penguin.AccountTypeDefault,
			},
		}
	} else if v := head.GetC2CCmd(); v != 0 {
		// private
		fromUin := int64(head.GetFromUin())
		from, ok := d.cntm.GetContact(uin, fromUin)
		if !ok {
			account, ok := d.accm.GetAccount(fromUin)
			if !ok {
				account = &penguin.Account{
					ID:   fromUin,
					Type: penguin.AccountTypeDefault,
				}
				_, _ = d.accm.SetAccount(account.ID, account)
			}
			from = &penguin.Contact{
				Account: account,
				Display: string(head.GetFromNick()),
			}
			_, _ = d.cntm.SetContact(uin, account.ID, from)
		}
		msg.From = &penguin.User{
			Account: from.Account,
			Display: from.Display,
		}
		toUin := int64(head.GetToUin())
		to, ok := d.cntm.GetContact(uin, toUin)
		if !ok {
			account, ok := d.accm.GetAccount(toUin)
			if !ok {
				account = &penguin.Account{
					ID:   toUin,
					Type: penguin.AccountTypeDefault,
				}
				_, _ = d.accm.SetAccount(account.ID, account)
			}
			to = &penguin.Contact{
				Account: account,
				Display: string(head.GetFromNick()),
			}
			_, _ = d.cntm.SetContact(uin, account.ID, to)
		}
		msg.Chat = &penguin.Chat{
			ID:   0,
			Type: penguin.ChatTypePrivate,
			User: &penguin.User{
				Account: to.Account,
				Display: to.Display,
			},
		}
	}
	_ = pgn.NewMessageEncoder(body).Encode(&msg)
	ph, _ := json.Marshal(head)
	pb, _ := json.Marshal(body)
	pm, _ := json.Marshal(msg)
	log.Debug("head:%s body:%s msg:%s", ph, pb, pm)
	switch msg.Chat.Type {
	case penguin.ChatTypePrivate:
		log.Chat("private:%d(%s) user:%d(%s) text:%s", msg.Chat.User.Account.ID, msg.Chat.User.Account.Username, msg.From.Account.ID, msg.From.Account.Username, msg.Text)
	case penguin.ChatTypeGroup:
		log.Chat("group:%d(%s) user:%d(%s) text:%s", msg.Chat.ID, msg.Chat.Title, msg.From.Account.ID, msg.From.Account.Username, msg.Text)
	}
	return nil
}

func (d *Daemon) OnSendMessage(msg *penguin.Message) error {
	return nil
}
