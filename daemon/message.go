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

func (d *Daemon) OnRecvMessage(head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) error {
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
	if chat := head.GetDiscussInfo(); chat != nil {
		// discuss
		msg.Chat = &penguin.Chat{
			ID:      int64(chat.GetDiscussUin()),
			Type:    penguin.ChatTypeDiscuss,
			Title:   string(chat.GetDiscussName()),
			Display: string(chat.GetDiscussRemark()),
		}
	} else if chat := head.GetDiscussInfo(); chat != nil {
		// discuss private
		msg.Chat = &penguin.Chat{
			ID:      int64(chat.GetDiscussUin()),
			Type:    penguin.ChatTypeDiscussPrivate,
			Title:   string(chat.GetDiscussName()),
			Display: string(chat.GetDiscussRemark()),
		}
		msg.Chat.User = &penguin.User{
			Account: &penguin.Account{
				ID:   int64(head.GetFromUin()),
				Type: penguin.AccountTypeDefault,
			},
		}
	} else if chat := head.GetGroupInfo(); chat != nil {
		// group
		msg.Chat = &penguin.Chat{
			ID:    int64(chat.GetGroupCode()),
			Type:  penguin.ChatTypeGroup,
			Title: string(chat.GetGroupName()),
		}
		msg.From.Display = string(chat.GetGroupCard())
	} else if chat := head.GetC2CTmpMsgHead(); chat != nil {
		// group private
		msg.Chat = &penguin.Chat{
			ID:   int64(chat.GetGroupCode()),
			Type: penguin.ChatTypeGroupPrivate,
		}
		msg.Chat.User = &penguin.User{
			Account: &penguin.Account{
				ID:   int64(head.GetFromUin()),
				Type: penguin.AccountTypeDefault,
			},
		}
	} else if cmd := head.GetC2CCmd(); cmd != 0 {
		// private
		msg.Chat = &penguin.Chat{
			ID:   0,
			Type: penguin.ChatTypePrivate,
		}
		msg.Chat.User = &penguin.User{
			Account: &penguin.Account{
				ID:   int64(head.GetFromUin()),
				Type: penguin.AccountTypeDefault,
			},
		}
		msg.From.Display = string(head.GetFromNick())
	}
	_ = pgn.NewMessageEncoder(body).Encode(&msg)
	ph, _ := json.Marshal(head)
	pb, _ := json.Marshal(body)
	p, _ := json.MarshalIndent(msg, "", "  ")
	log.Debug("head:%s body:%s\n%s", ph, pb, p)
	return nil
}

func (d *Daemon) OnSendMessage(msg *penguin.Message) error {
	return nil
}
