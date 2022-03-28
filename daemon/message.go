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
	"github.com/elap5e/penguin/pkg/log"
)

func (d *Daemon) OnRecvMessage(head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) error {
	msg := penguin.Message{
		MessageID: int64(head.GetMsgSeq()),
		From:      &penguin.User{Account: &penguin.Account{ID: int64(head.GetFromUin())}},
		Time:      int64(head.GetMsgTime()),
	}
	if chat := head.GetGroupInfo(); chat != nil {
		msg.Chat = &penguin.Chat{
			ID:    int64(chat.GetGroupCode()),
			Type:  "group",
			Title: string(chat.GetGroupName()),
		}
		msg.From.Display = string(chat.GetGroupCard())
	}
	msg.Text, msg.Entities, _ = encodeContent(head, body)
	ph, _ := json.MarshalIndent(head, "", "  ")
	pb, _ := json.MarshalIndent(body, "", "  ")
	p, _ := json.MarshalIndent(msg, "", "  ")
	log.Debug("OnRecvMessage:\n%s\n%s\n%s", ph, pb, p)
	return nil
}

func (d *Daemon) OnSendMessage(msg *penguin.Message) error {
	return nil
}

func encodeContent(head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) (string, []*penguin.MessageEntity, error) {
	text, entities := "", []*penguin.MessageEntity{}
	elems := body.GetRichText().GetElems()
	for _, elem := range elems {
		if elem := elem.GetText(); elem != nil {
			text += string(elem.GetStr())
		}
	}
	return text, entities, nil
}
