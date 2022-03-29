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

package pgn

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/message/face"
	"github.com/elap5e/penguin/daemon/message/pb"
)

type messageDecoder struct {
	msg *penguin.Message
}

func NewMessageDecoder(msg *penguin.Message) *messageDecoder {
	return &messageDecoder{
		msg: msg,
	}
}

func (me *messageDecoder) Decode(body *pb.IMMsgBody_MsgBody) error {
	return nil
}

type messageEncoder struct {
	body   *pb.IMMsgBody_MsgBody
	next   int
	buffer *bytes.Buffer
	offset int64
}

func NewMessageEncoder(body *pb.IMMsgBody_MsgBody) *messageEncoder {
	return &messageEncoder{
		body:   body,
		next:   0,
		buffer: bytes.NewBuffer(nil),
		offset: 0,
	}
}

func (me *messageEncoder) Encode(msg *penguin.Message) error {
	me.buffer.Reset()
	elems := me.body.GetRichText().GetElems()
	for me.next = 0; me.next < len(elems); me.next++ {
		elem := elems[me.next]
		if v := elem.GetText(); v != nil {
			me.encodeText(v, msg)
		} else if v := elem.GetFace(); v != nil {
			me.encodeFace(v, msg)
		} else if v := elem.GetNotOnlineImage(); v != nil {
			me.encodeNotOnlineImage(v, msg)
		} else if v := elem.GetMarketFace(); v != nil {
			me.encodeMarketFace(v, msg)
		} else if v := elem.GetShakeWindow(); v != nil {
			me.encodeShakeWindow(v, msg)
		} else if v := elem.GetCustomFace(); v != nil {
			me.encodeCustomFace(v, msg)
		} else if v := elem.GetCommonElem(); v != nil {
			me.encodeCommonElem(v, msg)
		}
	}
	msg.Text = me.buffer.String()
	return nil
}

func (me *messageEncoder) encodeText(elem *pb.IMMsgBody_Text, msg *penguin.Message) {
	buf := elem.GetAttr_6Buf()
	if len(buf) < 13 {
		n, _ := me.buffer.Write(elem.GetStr())
		me.offset += int64(n)
	} else {
		id := int64(buf[7])<<24 + int64(buf[8])<<16 + int64(buf[9])<<8 + int64(buf[10])
		n, _ := me.buffer.Write(elem.GetStr())
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
			Type:   "mention",
			Offset: me.offset,
			Length: int64(n),
			URL:    fmt.Sprintf("?id=%d", id),
			User:   &penguin.User{Account: &penguin.Account{ID: id}},
		})
		me.offset += int64(n)
	}
}

func (me *messageEncoder) encodeFace(elem *pb.IMMsgBody_Face, msg *penguin.Message) {
	id := elem.GetIndex()
	n, _ := me.buffer.Write([]byte(face.FaceType(id).String()))
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "face",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?id=%d", id),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeNotOnlineImage(elem *pb.IMMsgBody_NotOnlineImage, msg *penguin.Message) {
	n, _ := me.buffer.Write([]byte("[图片]"))
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x", elem.GetPicMd5()),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeCustomFace(elem *pb.IMMsgBody_CustomFace, msg *penguin.Message) {
	n, _ := me.buffer.Write([]byte("[图片]"))
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x", elem.GetMd5()),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeCommonElem(elem *pb.IMMsgBody_CommonElem, msg *penguin.Message) {
	switch elem.GetServiceType() {
	case 2: // poke
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
			Type:   "poke",
			Offset: 0,
			Length: 0,
			URL:    fmt.Sprintf("?id=%d", elem.GetBusinessType()),
		})
	case 33: // extra face
		info := pb.HummerCommonElement_MsgElemInfoServtype33{}
		_ = proto.Unmarshal(elem.GetPbElem(), &info)
		id := info.GetIndex()
		n, _ := me.buffer.Write([]byte(face.FaceType(id).String()))
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
			Type:   "face",
			Offset: me.offset,
			Length: int64(n),
			URL:    fmt.Sprintf("?id=%d", id),
		})
		me.offset += int64(n)
	case 37: // big face
		info := pb.HummerCommonElement_MsgElemInfoServtype37{}
		_ = proto.Unmarshal(elem.GetPbElem(), &info)
		me.next++
		name := me.body.GetRichText().GetElems()[me.next].GetText().GetStr()
		n, _ := me.buffer.Write(name)
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
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

func (me *messageEncoder) encodeShakeWindow(elem *pb.IMMsgBody_ShakeWindow, msg *penguin.Message) {
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "poke",
		Offset: 0,
		Length: 0,
		URL:    "?id=0",
	})
}

func (me *messageEncoder) encodeMarketFace(elem *pb.IMMsgBody_MarketFace, msg *penguin.Message) {
	me.next++
	name := elem.GetFaceName()
	if len(name) == 0 {
		name = me.body.GetRichText().GetElems()[me.next].GetText().GetStr()
	}
	n, _ := me.buffer.Write(name)
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
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
