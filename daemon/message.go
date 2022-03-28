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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/message/face"
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
	_ = NewMessageEncoder(head, body, &msg).Encode()
	ph, _ := json.Marshal(head)
	pb, _ := json.Marshal(body)
	p, _ := json.MarshalIndent(msg, "", "  ")
	log.Debug("head:%s body:%s\n%s", ph, pb, p)
	return nil
}

func (d *Daemon) OnSendMessage(msg *penguin.Message) error {
	return nil
}

type messageEncoder struct {
	head   *pb.MsgCommon_MsgHead
	body   *pb.IMMsgBody_MsgBody
	next   int
	buffer *bytes.Buffer
	offset int64
	msg    *penguin.Message
}

func NewMessageEncoder(head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody, msg *penguin.Message) *messageEncoder {
	return &messageEncoder{
		head:   head,
		body:   body,
		next:   0,
		buffer: bytes.NewBuffer(nil),
		offset: 0,
		msg:    msg,
	}
}

func (me *messageEncoder) Encode() error {
	elems := me.body.GetRichText().GetElems()
	for me.next = 0; me.next < len(elems); me.next++ {
		elem := elems[me.next]
		if v := elem.GetText(); v != nil {
			me.encodeText(v)
		} else if v := elem.GetFace(); v != nil {
			me.encodeFace(v)
		} else if v := elem.GetNotOnlineImage(); v != nil {
			me.encodeNotOnlineImage(v)
		} else if v := elem.GetMarketFace(); v != nil {
			me.encodeMarketFace(v)
		} else if v := elem.GetShakeWindow(); v != nil {
			me.encodeShakeWindow(v)
		} else if v := elem.GetCustomFace(); v != nil {
			me.encodeCustomFace(v)
		} else if v := elem.GetCommonElem(); v != nil {
			me.encodeCommonElem(v)
		}
	}
	me.msg.Text = me.buffer.String()
	return nil
}

func (me *messageEncoder) encodeText(elem *pb.IMMsgBody_Text) {
	buf := elem.GetAttr_6Buf()
	if len(buf) < 13 {
		n, _ := me.buffer.Write(elem.GetStr())
		me.offset += int64(n)
	} else {
		id := int64(buf[7])<<24 + int64(buf[8])<<16 + int64(buf[9])<<8 + int64(buf[10])
		n, _ := me.buffer.Write(elem.GetStr())
		me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
			Type:   "mention",
			Offset: me.offset,
			Length: int64(n),
			URL:    fmt.Sprintf("?id=%d", id),
			User:   &penguin.User{Account: &penguin.Account{ID: id}},
		})
		me.offset += int64(n)
	}
}

func (me *messageEncoder) encodeFace(elem *pb.IMMsgBody_Face) {
	id := elem.GetIndex()
	n, _ := me.buffer.Write([]byte(face.FaceType(id).String()))
	me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
		Type:   "face",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?id=%d", id),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeNotOnlineImage(elem *pb.IMMsgBody_NotOnlineImage) {
	n, _ := me.buffer.Write([]byte("[图片]"))
	me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x", elem.GetPicMd5()),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeCustomFace(elem *pb.IMMsgBody_CustomFace) {
	n, _ := me.buffer.Write([]byte("[图片]"))
	me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x", elem.GetMd5()),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeCommonElem(elem *pb.IMMsgBody_CommonElem) {
	switch elem.GetServiceType() {
	case 2: // poke
		me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
			Type:   "poke",
			Offset: 0,
			Length: 0,
			URL:    fmt.Sprintf("?id=%d", elem.GetBusinessType()),
		})
	case 33: // extra face
		info := pb.HummerCommelem_MsgElemInfoServtype33{}
		_ = proto.Unmarshal(elem.GetPbElem(), &info)
		id := info.GetIndex()
		n, _ := me.buffer.Write([]byte(face.FaceType(id).String()))
		me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
			Type:   "face",
			Offset: me.offset,
			Length: int64(n),
			URL:    fmt.Sprintf("?id=%d", id),
		})
		me.offset += int64(n)
	case 37: // big face
		info := pb.HummerCommelem_MsgElemInfoServtype37{}
		_ = proto.Unmarshal(elem.GetPbElem(), &info)
		me.next++
		name := me.body.GetRichText().GetElems()[me.next].GetText().GetStr()
		n, _ := me.buffer.Write(name)
		me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
			Type:   "face",
			Offset: me.offset,
			Length: int64(n),
			URL: fmt.Sprintf(
				"?id=%d&pid=%s&sid=%s",
				info.GetQsid(),
				base64.RawURLEncoding.EncodeToString(info.GetPackid()),
				base64.RawURLEncoding.EncodeToString(info.GetStickerid()),
			),
		})
		me.offset += int64(n)
	}
}

func (me *messageEncoder) encodeShakeWindow(elem *pb.IMMsgBody_ShakeWindow) {
	me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
		Type:   "poke",
		Offset: 0,
		Length: 0,
		URL:    "?id=0",
	})
}

func (me *messageEncoder) encodeMarketFace(elem *pb.IMMsgBody_MarketFace) {
	me.next++
	name := elem.GetFaceName()
	if len(name) == 0 {
		name = me.body.GetRichText().GetElems()[me.next].GetText().GetStr()
	}
	n, _ := me.buffer.Write(name)
	me.msg.Entities = append(me.msg.Entities, &penguin.MessageEntity{
		Type:   "sticker",
		Offset: me.offset,
		Length: int64(n),
		URL: fmt.Sprintf(
			"?id=%s&tid=%d&key=%s&h=%d&w+%d",
			base64.RawURLEncoding.EncodeToString(elem.GetFaceId()),
			elem.GetTabId(),
			base64.RawURLEncoding.EncodeToString(elem.GetKey()),
			elem.GetImageHeight(),
			elem.GetImageWidth(),
		),
	})
	me.offset += int64(n)
}
