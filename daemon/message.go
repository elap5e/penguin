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
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/encoding/pgn"
	"github.com/elap5e/penguin/pkg/log"
)

// mask 0x00000000ffffffff
func (d *Daemon) prefetchAccount(id int64) error {
	account, ok := d.accm.GetAccount(id)
	if !ok {
		account = &penguin.Account{
			ID:   id,
			Type: penguin.AccountTypeDefault,
		}
		_, _ = d.accm.SetAccount(account.ID, account)
	}
	return nil
}

func (d *Daemon) getOrLoadContact(id, cid int64) *penguin.Contact {
	contact, ok := d.cntm.GetContact(id, cid)
	if !ok {
		account, _ := d.accm.GetAccount(cid)
		contact = &penguin.Contact{
			Account: account,
		}
		_, _ = d.cntm.SetContact(id, account.ID, contact)
	}
	return contact
}

func (d *Daemon) getOrLoadChatUser(cid, uid int64) *penguin.User {
	from, ok := d.chtm.GetUser(cid, uid)
	if !ok {
		account, _ := d.accm.GetAccount(uid)
		from = &penguin.User{
			Account: account,
		}
		_, _ = d.chtm.SetUser(cid, account.ID, from)
	}
	return from
}

func (d *Daemon) getOrLoadChatGroup(id int64, name string) *penguin.Chat {
	chat, ok := d.chtm.GetChat(id)
	if !ok {
		chat = &penguin.Chat{
			ID:    id,
			Type:  penguin.ChatTypeGroup,
			Title: name,
		}
		_, _ = d.chtm.SetChat(id, chat)
	}
	return chat
}

func (d *Daemon) getOrLoadChatGroupPrivate(id int64, name string) *penguin.Chat {
	chat, ok := d.chtm.GetChat(id)
	if !ok {
		chat = &penguin.Chat{
			ID:      id,
			Type:    penguin.ChatTypeGroup,
			Display: name,
		}
	}
	return &penguin.Chat{
		ID:      chat.ID,
		Type:    penguin.ChatTypeGroupPrivate,
		Display: chat.Display,
	}
}

func (d *Daemon) OnRecvMessage(id int64, head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) error {
	msg := penguin.Message{
		MessageID: int64(head.GetMsgSeq()),
		Time:      int64(head.GetMsgTime()),
	}
	// pre-fetch accounts
	_ = d.prefetchAccount(int64(head.GetFromUin()))
	_ = d.prefetchAccount(int64(head.GetToUin()))
	// identify message type
	if v := head.GetDiscussInfo(); v != nil {
		// discuss
	} else if v := head.GetDiscussInfo(); v != nil {
		// discuss private
	} else if v := head.GetGroupInfo(); v != nil {
		// group
		gid, fid := int64(v.GetGroupCode()), int64(head.GetFromUin())
		msg.Chat = d.getOrLoadChatGroup(gid, string(v.GetGroupName()))
		msg.From = d.getOrLoadChatUser(gid, fid)
		_, _ = d.chtm.SetChatSeq(0, gid, 0, head.GetMsgSeq())
	} else if v := head.GetC2CTmpMsgHead(); v != nil {
		// group private
		gid, fid, tid := int64(v.GetGroupCode()), int64(head.GetFromUin()), int64(head.GetToUin())
		msg.Chat = d.getOrLoadChatGroupPrivate(gid, "")
		msg.From = d.getOrLoadChatUser(gid, fid)
		// check if the sender is self
		if id == fid && id != tid {
			msg.Chat.User = d.getOrLoadChatUser(gid, tid)
		} else {
			msg.Chat.User = msg.From
			// tid = fid
		}
		// _, _ = d.chtm.SetChatSeq(id, gid, tid, head.GetMsgSeq())
	} else if v := head.GetC2CCmd(); v != 0 {
		// private
		msg.Chat = &penguin.Chat{Type: penguin.ChatTypePrivate}
		fid, tid := int64(head.GetFromUin()), int64(head.GetToUin())
		from := d.getOrLoadContact(id, fid)
		msg.From = &penguin.User{
			Account: from.Account,
			Display: from.Display,
		}
		// check if the sender is self
		if id == fid && id != tid {
			to := d.getOrLoadContact(id, tid)
			msg.Chat.User = &penguin.User{
				Account: to.Account,
				Display: to.Display,
			}
		} else {
			msg.Chat.User = msg.From
			// tid = fid
		}
		// _, _ = d.chtm.SetChatSeq(id, 0, tid, head.GetMsgSeq())
	}
	if err := pgn.NewMessageEncoder(body).Encode(&msg); err != nil {
		return err
	}
	return d.onRecvMessage(id, head, body, &msg)
}

func (d *Daemon) SendMessage(id int64, msg *penguin.Message) error {
	var req pb.MsgService_PbSendMsgReq
	req.RoutingHead = &pb.MsgService_RoutingHead{}
	// identify message type
	if msg.Chat.Type == penguin.ChatTypeDiscuss {
		// discuss
	} else if msg.Chat.Type == penguin.ChatTypeDiscussPrivate {
		// discuss private
	} else if msg.Chat.Type == penguin.ChatTypeGroup {
		// group
		req.RoutingHead.Grp = &pb.MsgService_Grp{
			GroupCode: uint64(msg.Chat.ID),
		}
		req.MsgSeq, _ = d.chtm.GetNextChatSeq(0, msg.Chat.ID, 0)
	} else if msg.Chat.Type == penguin.ChatTypeGroupPrivate {
		// group private
		req.RoutingHead.GrpTmp = &pb.MsgService_GrpTmp{
			GroupUin: uint64(msg.Chat.ID),
			ToUin:    uint64(msg.Chat.User.Account.ID),
		}
		req.MsgSeq, _ = d.chtm.GetNextChatSeq(id, msg.Chat.ID, msg.Chat.User.Account.ID)
	} else if msg.Chat.Type == penguin.ChatTypePrivate {
		// private
		req.RoutingHead.C2C = &pb.MsgService_C2C{
			ToUin: uint64(msg.Chat.User.Account.ID),
		}
		req.MsgSeq, _ = d.chtm.GetNextChatSeq(id, 0, msg.Chat.User.Account.ID)
	} else {
		return d.SendChannelMessage(id, msg)
	}
	// encode message
	req.ContentHead = &pb.MsgCommon_ContentHead{}
	req.MsgBody = &pb.IMMsgBody_MsgBody{}
	if err := pgn.NewMessageDecoder(req.MsgBody).Decode(msg); err != nil {
		return err
	}
	resp, err := d.msgm.SendMessage(id, &req)
	if err != nil {
		return err
	}
	return d.onSendMessage(id, &req, resp, msg)
}

func (d *Daemon) onRecvMessage(id int64, head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody, msg *penguin.Message) error {
	go d.pushMessage(msg)
	go d.fetchBlobs(msg)
	ph, _ := json.Marshal(head)
	pb, _ := json.Marshal(body)
	pm, _ := json.Marshal(msg)
	log.Debug("id:%d head:%s body:%s msg:%s", id, ph, pb, pm)
	log.Chat(id, msg)
	return nil
}

func (d *Daemon) onSendMessage(id int64, req *pb.MsgService_PbSendMsgReq, resp *pb.MsgService_PbSendMsgResp, msg *penguin.Message) error {
	go d.pushMessage(msg)
	preq, _ := json.Marshal(req)
	prsp, _ := json.Marshal(resp)
	pmsg, _ := json.Marshal(msg)
	log.Debug("id:%d req:%s resp:%s msg:%s", id, preq, prsp, pmsg)
	log.Chat(id, msg)
	return nil
}

func (d *Daemon) pushMessage(msg *penguin.Message) {
	d.msgChan <- msg
}

func (d *Daemon) fetchBlobs(msg *penguin.Message) error {
	for _, v := range msg.Entities {
		switch v.Type {
		case "photo":
			if err := d.fetchBlob(v.Type, v.URL); err != nil {
				log.Error("fetchBlob %s%s failed, error:%v", v.Type, v.URL, err)
			}
		}
	}
	return nil
}

func (d *Daemon) fetchBlob(typ, str string) error {
	u, _ := url.Parse(str)
	query := u.Query()
	homepath, _ := os.UserHomeDir()
	basepath := path.Join(homepath, ".penguin", "cache")
	filepath := path.Join(basepath, "blobs", "md5", query.Get("md5"))
	if _, err := os.Stat(filepath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	url := ""
	switch typ {
	case "photo":
		url = fmt.Sprintf("https://gchat.qpic.cn/gchatpic_new/0/0-0-%s/0?term=2", strings.ToUpper(query.Get("md5")))
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code:%d", resp.StatusCode)
	}
	hash := sha256.New()
	body := io.TeeReader(resp.Body, hash)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}
	hashpath := fmt.Sprintf("sha256/%x", hash.Sum(nil)[:sha256.Size])
	if err := os.Rename(filepath, path.Join(basepath, "blobs", hashpath)); err != nil {
		return err
	}
	if err := os.Symlink("../"+hashpath, filepath); err != nil {
		return err
	}
	head := make([]byte, 512)
	_, _ = file.ReadAt(head, 0)
	switch typ {
	case "audio":
		filepath = path.Join(basepath, "audio", query.Get("md5"))
	case "photo":
		filepath = path.Join(basepath, "photo", query.Get("md5"))
	case "video":
		filepath = path.Join(basepath, "video", query.Get("md5"))
	case "voice":
		filepath = path.Join(basepath, "voice", query.Get("md5"))
	}
	switch http.DetectContentType(head) {
	default:
		return nil
	case "image/bmp":
		filepath += ".bmp"
	case "image/gif":
		filepath += ".gif"
	case "image/jpeg":
		filepath += ".jpeg"
	case "image/png":
		filepath += ".png"
	case "image/webp":
		filepath += ".webp"
	case "image/x-icon":
		filepath += ".ico"
	}
	return os.Symlink("../blobs/"+hashpath, filepath)
}
