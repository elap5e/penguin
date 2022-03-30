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

package log

import (
	"fmt"

	"github.com/elap5e/penguin"
)

func Chat(id int64, msg *penguin.Message) {
	str := fmt.Sprintf("[%d] message not parsed", id)
	switch msg.Chat.Type {
	case penguin.ChatTypePrivate:
		str = fmt.Sprintf(
			"[%d] private:%d(%s) user:%d(%s) text:%s",
			id,
			msg.Chat.User.Account.ID,
			msg.Chat.User.Account.Username,
			msg.From.Account.ID,
			msg.From.Account.Username,
			msg.Text,
		)
	case penguin.ChatTypeGroup:
		str = fmt.Sprintf(
			"[%d] group:%d(%s) user:%d(%s) text:%s",
			id,
			msg.Chat.ID,
			msg.Chat.Title,
			msg.From.Account.ID,
			msg.From.Account.Username,
			msg.Text,
		)
	case penguin.ChatTypeGroupPrivate:
		str = fmt.Sprintf(
			"[%d] group:%d(%s) private:%d(%s) user:%d(%s) text:%s",
			id,
			msg.Chat.ID,
			msg.Chat.Title,
			msg.Chat.User.Account.ID,
			msg.Chat.User.Account.Username,
			msg.From.Account.ID,
			msg.From.Account.Username,
			msg.Text,
		)
	case penguin.ChatTypeRoomText:
		str = fmt.Sprintf(
			"[%d] channel:0x%x(%s) room:0x%x(%s) user:0x%x(%s) text:%s",
			id,
			msg.Chat.Chat.ID,
			msg.Chat.Chat.Title,
			msg.Chat.ID,
			msg.Chat.Title,
			msg.From.Account.ID,
			msg.From.Account.Username,
			msg.Text,
		)
	}
	logger.Println("\x1b[37;1m[CHAT]", str+"\x1b[0m")
}
