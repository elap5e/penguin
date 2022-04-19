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
	chatDisplay := msg.Chat.Display
	if chatDisplay == "" {
		chatDisplay = msg.Chat.Title
	}
	fromDisplay := msg.From.Display
	if fromDisplay == "" {
		fromDisplay = msg.From.Account.Username
	}
	switch msg.Chat.Type {
	case penguin.ChatTypeGroup:
		str = fmt.Sprintf(
			"[%d] group:%d(%s) user:%d(%s) text:%q",
			id,
			msg.Chat.ID,
			chatDisplay,
			msg.From.Account.ID,
			fromDisplay,
			msg.Text,
		)
	case penguin.ChatTypeGroupPrivate:
		userDisplay := msg.Chat.User.Display
		if userDisplay == "" {
			userDisplay = msg.Chat.User.Account.Username
		}
		str = fmt.Sprintf(
			"[%d] group:%d(%s) private:%d(%s) user:%d(%s) text:%q",
			id,
			msg.Chat.ID,
			chatDisplay,
			msg.Chat.User.Account.ID,
			userDisplay,
			msg.From.Account.ID,
			fromDisplay,
			msg.Text,
		)
	case penguin.ChatTypePrivate:
		userDisplay := msg.Chat.User.Display
		if userDisplay == "" {
			userDisplay = msg.Chat.User.Account.Username
		}
		str = fmt.Sprintf(
			"[%d] private:%d(%s) user:%d(%s) text:%q",
			id,
			msg.Chat.User.Account.ID,
			userDisplay,
			msg.From.Account.ID,
			fromDisplay,
			msg.Text,
		)
	case penguin.ChatTypeRoomText, penguin.ChatTypeRoomLive:
		str = fmt.Sprintf(
			"[%d] channel:%d(%s) room:%d(%s) user:%d(%s) text:%q",
			id,
			msg.Chat.Channel.ID,
			msg.Chat.Channel.Title,
			msg.Chat.ID,
			chatDisplay,
			msg.From.Account.ID,
			fromDisplay,
			msg.Text,
		)
	case penguin.ChatTypeRoomPrivate:
		str = fmt.Sprintf(
			"[%d] channel:%d(%s) private:%d:%d user:%d(%s) text:%q",
			id,
			msg.Chat.Channel.Channel.ID,
			msg.Chat.Channel.Channel.Title,
			msg.Chat.Channel.ID,
			msg.Chat.ID,
			msg.From.Account.ID,
			fromDisplay,
			msg.Text,
		)
	}
	logger.Println("\x1b[37;1m[CHAT]", str+"\x1b[0m")
}
