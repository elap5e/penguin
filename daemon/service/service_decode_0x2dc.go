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

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/message/dto"
	"github.com/elap5e/penguin/daemon/service/pb"
	"github.com/elap5e/penguin/pkg/log"
)

func (m *Manager) Decode0x2dc(uin int64, msg *dto.Message) error {
	if len(msg.MessageBytes) < 5 {
		return fmt.Errorf("invalid message length: %d < 5", len(msg.MessageBytes))
	}
	chatID := int64(binary.BigEndian.Uint32(msg.MessageBytes[:4]))
	chat, _ := m.GetChatManager().GetChat(chatID)

	typ := msg.MessageBytes[4]
	switch typ {
	case 0x0c:
		// chat mute
		return m.onChatMute(uin, chat, msg)
	case 0x0e:
		// chat anonymous
		return m.onChatAnonymous(uin, chat, msg)
	case 0x10, 0x11, 0x14, 0x15:
		// chat tips
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
			chat, _ = m.GetChatManager().GetChat(chatID)
		}
		if v := notify.GetMsgGraytips(); v != nil {
			return m.onChatMessageGrayTips(uin, chat, v)
		} else if v := notify.GetMsgGroupNotify(); v != nil {
			return m.onChatMessageNotify(uin, chat, v)
		} else if v := notify.GetMsgRecall(); v != nil {
			return m.onChatMessageRecall(uin, chat, v)
		} else if v := notify.GetGeneralGrayTip(); v != nil {
			return m.onChatGrayTips(uin, chat, v)
		}
	case 0x03:
		fallthrough
	case 0x0f:
		fallthrough
	default:
		dumpUnknown(msg.Type, msg)
	}
	return nil
}

func (m *Manager) onChatMute(uin int64, chat *penguin.Chat, msg *dto.Message) error {
	fromID := int64(binary.BigEndian.Uint32(msg.MessageBytes[6:10]))
	fromDisplay := m.getChatUserDisplay(chat.ID, fromID)
	// time := binary.BigEndian.Uint32(msg.MessageBytes[10:14])
	for i := 0; i < int(binary.BigEndian.Uint16(msg.MessageBytes[14:16])); i++ {
		toID := int64(binary.BigEndian.Uint32(msg.MessageBytes[16+i*8 : 20+i*8]))
		seconds := time.Second * time.Duration(binary.BigEndian.Uint32(msg.MessageBytes[20+i*8:24+i*8]))
		if toID == 0 {
			if seconds == 0 {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]解除全员禁言", uin, chat.ID, chat.Title, fromID, fromDisplay)
			} else {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]全员禁言", uin, chat.ID, chat.Title, fromID, fromDisplay)
			}
		} else {
			toDisplay := m.getChatUserDisplay(chat.ID, toID)
			if seconds == 0 {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]解除[%d(%s)]禁言", uin, chat.ID, chat.Title, fromID, fromDisplay, toID, toDisplay)
			} else {
				log.Notifyf("[%d] group:%d(%s) [%d(%s)]禁言[%d(%s)][%s]", uin, chat.ID, chat.Title, fromID, fromDisplay, toID, toDisplay, seconds)
			}
		}
	}
	return nil
}

func (m *Manager) onChatMessageGrayTips(uin int64, chat *penguin.Chat, notify *pb.ChatTips_AIOGrayTipsInfo) error {
	return nil
}

func (m *Manager) onChatMessageNotify(uin int64, chat *penguin.Chat, notify *pb.ChatTips_GroupNotifyInfo) error {
	return nil
}

func (m *Manager) onChatMessageRecall(uin int64, chat *penguin.Chat, notify *pb.ChatTips_MessageRecallReminder) error {
	userID, suffix := int64(notify.GetUin()), notify.GetMsgWordingInfo().GetItemName()
	userDisplay := m.getChatUserDisplay(chat.ID, userID)
	if suffix != "" && !strings.HasPrefix(suffix, "，") {
		suffix = "，" + suffix
	}
	for _, v := range notify.GetRecalledMsgList() {
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]撤回了一条消息[%d:%d]%s", uin, chat.ID, chat.Title, userID, userDisplay, v.GetTime(), v.GetSeq(), suffix)
	}
	return nil
}

func (m *Manager) onChatGrayTips(uin int64, chat *penguin.Chat, notify *pb.ChatTips_GeneralGrayTipInfo) error {
	id, tmplID := notify.GetBusiId(), notify.GetTemplId()
	if id == 1003 {
		userID, userDisplay, lastID, lastDisplay, honor := int64(0), "", int64(0), "", ""
		for _, v := range notify.GetMsgTemplParam() {
			if bytes.Equal(v.GetName(), []byte("uin")) {
				userID, _ = strconv.ParseInt(string(v.GetValue()), 10, 64)
			} else if bytes.Equal(v.GetName(), []byte("nick")) {
				userDisplay = string(v.GetValue())
			} else if bytes.Equal(v.GetName(), []byte("uin_last")) {
				lastID, _ = strconv.ParseInt(string(v.GetValue()), 10, 64)
			} else if bytes.Equal(v.GetName(), []byte("nick_last")) {
				lastDisplay = string(v.GetValue())
			} else if bytes.Equal(v.GetName(), []byte("honor_name_2")) {
				honor = string(v.GetValue())
			}
		}
		if userDisplay == "" {
			userDisplay = m.getChatUserDisplay(chat.ID, userID)
		}
		if lastDisplay == "" {
			lastDisplay = m.getChatUserDisplay(chat.ID, lastID)
		}
		if honor == "" && (tmplID == 1053 || tmplID == 1054 || tmplID == 10094) {
			honor = "龙王"
		}
		if tmplID == 1052 {
			log.Notifyf("[%d] group:%d(%s) [%d(%s)]在群聊中连续发消息超过7天, 获得[%s]标识。", uin, chat.ID, chat.Title, userID, userDisplay, honor)
		} else if tmplID == 1053 {
			log.Notifyf("[%d] group:%d(%s) 昨日[%d(%s)]在群内发言最积极，夺走了[%d(%s)]的[%s]标识。", uin, chat.ID, chat.Title, userID, userDisplay, lastID, lastDisplay, honor)
		} else if tmplID == 1054 {
			log.Notifyf("[%d] group:%d(%s) 昨日[%d(%s)]在群内发言最积极，获得[%s]标识。", uin, chat.ID, chat.Title, userID, userDisplay, honor)
		} else if tmplID == 1067 {
			log.Notifyf("[%d] group:%d(%s) [%d(%s)]在群聊中连续发表情包超过3天，且累计数量超过20条，获得[%s]标识。", uin, chat.ID, chat.Title, userID, userDisplay, honor)
		} else if tmplID == 10094 {
			log.Notifyf("[%d] group:%d(%s) 昨日你在群内发言最积极，夺走了[%d(%s)]的[%s]标识。", uin, chat.ID, chat.Title, lastID, lastDisplay, honor)
		}
	} else if id == 1061 {
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
		fromDisplay := m.getChatUserDisplay(chat.ID, fromID)
		toDisplay := m.getChatUserDisplay(chat.ID, toID)
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]%s[%d(%s)]%s", uin, chat.ID, chat.Title, fromID, fromDisplay, action, toID, toDisplay, suffix)
	} else if id == 1068 {
		userID, userDisplay, action := int64(0), "", ""
		for _, v := range notify.GetMsgTemplParam() {
			if bytes.Equal(v.GetName(), []byte("mqq_uin")) {
				userID, _ = strconv.ParseInt(string(v.GetValue()), 10, 64)
			} else if bytes.Equal(v.GetName(), []byte("mqq_nick")) && len(v.GetValue()) != 0 {
				userDisplay = string(v.GetValue())
			} else if bytes.Equal(v.GetName(), []byte("user_sign")) && len(v.GetValue()) != 0 {
				action = string(v.GetValue())
			}
		}
		if userDisplay == "" {
			userDisplay = m.getChatUserDisplay(chat.ID, userID)
		}
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]%s", uin, chat.ID, chat.Title, userID, userDisplay, action)
	}
	return nil
}

func (m *Manager) onChatAnonymous(uin int64, chat *penguin.Chat, msg *dto.Message) error {
	fromID := int64(binary.BigEndian.Uint32(msg.MessageBytes[6:10]))
	fromDisplay := m.getChatUserDisplay(chat.ID, fromID)
	if bytes.Equal(msg.MessageBytes[10:14], []byte{0, 0, 0, 0}) {
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]允许群内匿名聊天", uin, chat.ID, chat.Title, fromID, fromDisplay)
	} else {
		log.Notifyf("[%d] group:%d(%s) [%d(%s)]禁止群内匿名聊天", uin, chat.ID, chat.Title, fromID, fromDisplay)
	}
	return nil
}

func (m *Manager) getChatUserDisplay(chatID, userID int64) string {
	user, ok := m.GetChatManager().GetChatUser(chatID, userID)
	if ok {
		if user.Display != "" {
			return user.Display
		} else {
			return user.Account.Username
		}
	}
	return ""
}
