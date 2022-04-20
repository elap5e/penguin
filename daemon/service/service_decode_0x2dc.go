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

package service

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/message/dto"
	"github.com/elap5e/penguin/daemon/service/pb"
	"github.com/elap5e/penguin/pkg/log"
)

func (m *Manager) Decode0x2dc(uin int64, msg *dto.Message) error {
	if len(msg.MessageBytes) < 5 {
		return fmt.Errorf("invalid message length: %d < 5", len(msg.MessageBytes))
	}
	chatID := int64(binary.BigEndian.Uint32(msg.MessageBytes[:4]))
	typ := msg.MessageBytes[4]

	switch typ {
	case 0x0c: // ChatMute
		_ = m.onChatMute(uin, chatID, msg)
	case 0x10, 0x11, 0x14, 0x15: // ChatTips
		if len(msg.MessageBytes) < 8 {
			return fmt.Errorf("invalid message length: %d < 8", len(msg.MessageBytes))
		}
		var notify pb.ChatTips_NotifyMsgBody
		if err := proto.Unmarshal(msg.MessageBytes[7:], &notify); err != nil {
			return err
		}
		p, _ := json.Marshal(&notify)
		log.Debug("chat:%d dump notify:%s", chatID, string(p))
		if chatID != int64(notify.GetGroupCode()) {
			log.Warn("decode0x2dc chatID:%d != notify.GroupCode:%d", chatID, notify.GetGroupCode())
			chatID = int64(notify.GetGroupCode())
		}
		if v := notify.GetMsgGraytips(); v != nil {
			_ = m.onChatMessageGrayTips(uin, chatID, v)
		} else if v := notify.GetMsgGroupNotify(); v != nil {
			_ = m.onChatMessageNotify(uin, chatID, v)
		} else if v := notify.GetMsgRecall(); v != nil {
			_ = m.onChatMessageRecall(uin, chatID, v)
		} else if v := notify.GetGeneralGrayTip(); v != nil {
			_ = m.onChatGrayTips(uin, chatID, v)
		} else {
		}
	case 0x03:
		fallthrough
	case 0x0e:
		fallthrough
	case 0x0f:
		fallthrough
	default:
		dumpUnknown(msg.Type, msg)
	}

	return nil
}

func (m *Manager) onChatMute(uin, chatID int64, msg *dto.Message) error {
	fromID := int64(binary.BigEndian.Uint32(msg.MessageBytes[6:10]))
	chat, _ := m.GetChatManager().GetChat(chatID)
	from, _ := m.GetChatManager().GetChatUser(chatID, fromID)
	fromDisplay := from.Display
	if fromDisplay == "" {
		fromDisplay = from.Account.Username
	}
	// time := binary.BigEndian.Uint32(msg.MessageBytes[10:14])
	for i := 0; i < int(binary.BigEndian.Uint16(msg.MessageBytes[14:16])); i++ {
		toID := int64(binary.BigEndian.Uint32(msg.MessageBytes[16+i*8 : 20+i*8]))
		seconds := time.Second * time.Duration(binary.BigEndian.Uint32(msg.MessageBytes[20+i*8:24+i*8]))
		if toID == 0 {
			if seconds == 0 {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]解除全员禁言", uin, chat.ID, chat.Title, from.Account.ID, fromDisplay)
			} else {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]全员禁言", uin, chat.ID, chat.Title, from.Account.ID, fromDisplay)
			}
		} else {
			to, _ := m.GetChatManager().GetChatUser(chatID, toID)
			toDisplay := to.Display
			if toDisplay == "" {
				toDisplay = to.Account.Username
			}
			if seconds == 0 {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]解除[%d(%s)]禁言", uin, chat.ID, chat.Title, from.Account.ID, fromDisplay, to.Account.ID, toDisplay)
			} else {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]禁言[%d(%s)]%s", uin, chat.ID, chat.Title, from.Account.ID, fromDisplay, to.Account.ID, toDisplay, seconds)
			}
		}
	}
	return nil
}

func (m *Manager) onChatMessageGrayTips(uin, chatID int64, notify *pb.ChatTips_AIOGrayTipsInfo) error {
	return nil
}

func (m *Manager) onChatMessageNotify(uin, chatID int64, notify *pb.ChatTips_GroupNotifyInfo) error {
	return nil
}

func (m *Manager) onChatMessageRecall(uin, chatID int64, notify *pb.ChatTips_MessageRecallReminder) error {
	userID, suffix := int64(notify.GetUin()), notify.GetMsgWordingInfo().GetItemName()
	chat, _ := m.GetChatManager().GetChat(chatID)
	user, _ := m.GetChatManager().GetChatUser(chatID, userID)
	userDisplay := user.Display
	if userDisplay == "" {
		userDisplay = user.Account.Username
	}
	if suffix != "" && !strings.HasPrefix(suffix, "，") {
		suffix = "，" + suffix
	}
	for _, v := range notify.GetRecalledMsgList() {
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]撤回了一条消息[%d:%d]%s", uin, chat.ID, chat.Title, user.Account.ID, userDisplay, v.GetTime(), v.GetSeq(), suffix)
	}
	return nil
}

func (m *Manager) onChatGrayTips(uin, chatID int64, notify *pb.ChatTips_GeneralGrayTipInfo) error {
	id := notify.GetBusiId()
	if id == 1061 {
		fromID, action, toID, suffix := int64(0), "戳了戳", int64(0), ""
		for _, v := range notify.GetMsgTemplParam() {
			if bytes.Equal(v.GetName(), []byte("uin_str1")) {
				fromID, _ = strconv.ParseInt(string(v.GetValue()), 10, 64)
			} else if bytes.Equal(v.GetName(), []byte("action_str")) && len(v.GetValue()) != 0 {
				action = string(v.GetValue())
			} else if bytes.Equal(v.GetName(), []byte("uin_str2")) {
				toID, _ = strconv.ParseInt(string(v.GetValue()), 10, 64)
			} else if bytes.Equal(v.GetName(), []byte("suffix_str")) && len(v.GetValue()) != 0 {
				suffix = string(v.GetValue())
			}
		}
		chat, _ := m.GetChatManager().GetChat(chatID)
		from, _ := m.GetChatManager().GetChatUser(chatID, fromID)
		fromDisplay := from.Display
		if fromDisplay == "" {
			fromDisplay = from.Account.Username
		}
		to, _ := m.GetChatManager().GetChatUser(chatID, toID)
		toDisplay := to.Display
		if toDisplay == "" {
			toDisplay = to.Account.Username
		}
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]%s[%d(%s)]%s", uin, chat.ID, chat.Title, from.Account.ID, fromDisplay, action, to.Account.ID, toDisplay, suffix)
	} else if id == 1067 {
		userID, userDisplay, honor := int64(0), "", ""
		for _, v := range notify.GetMsgTemplParam() {
			if bytes.Equal(v.GetName(), []byte("uin")) {
				userID, _ = strconv.ParseInt(string(v.GetValue()), 10, 64)
			} else if bytes.Equal(v.GetName(), []byte("nick")) {
				userDisplay = string(v.GetValue())
			} else if bytes.Equal(v.GetName(), []byte("honor_name_2")) {
				honor = string(v.GetValue())
			}
		}
		chat, _ := m.GetChatManager().GetChat(chatID)
		user, _ := m.GetChatManager().GetChatUser(chatID, userID)
		if userDisplay == "" {
			userDisplay = user.Account.Username
		}
		// log.Notifyf("[%d] group:%d(%s) 昨日[%d(%s)]在群内发言最积极，获得[%s]标识。", uin, chat.ID, chat.Title, userID, userDisplay, honor)
		// log.Notifyf("[%d] group:%d(%s) [%d(%s)]在群聊中连续发消息超过7天, 获得[%s]标识。", uin, chat.ID, chat.Title, userID, userDisplay, honor)
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]在群聊中连续发表情包超过3天，且累计数量超过20条，获得[%s]标识。", uin, chat.ID, chat.Title, userID, userDisplay, honor)
	}
	return nil
}
