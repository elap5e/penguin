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
	"encoding/hex"
	"fmt"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/message/face"
	"github.com/elap5e/penguin/daemon/message/pb"
)

type messageDecoder struct {
	body   *pb.IMMsgBody_MsgBody
	next   int
	buffer *bytes.Buffer
	length int64
	offset int64
}

func NewMessageDecoder(body *pb.IMMsgBody_MsgBody) *messageDecoder {
	return &messageDecoder{
		body:   body,
		next:   0,
		buffer: nil,
		length: 0,
		offset: 0,
	}
}

func (md *messageDecoder) Decode(msg *penguin.Message) error {
	md.buffer = bytes.NewBuffer([]byte(msg.Text))
	md.length, md.offset = int64(md.buffer.Len()), 0
	elems := []*pb.IMMsgBody_Elem{}
	entities := msg.Entities
	for md.next = 0; md.next < len(entities); md.next++ {
		entity := entities[md.next]
		if md.offset < entity.Offset {
			md.decodeText(entity.Offset-md.offset, &elems)
		}
		if entity.Type == "face" {
			md.decodeFace(entity, &elems)
		} else if entity.Type == "photo" {
			md.decodePhoto(entity, &elems, msg)
		}
	}
	if md.offset < md.length {
		md.decodeText(md.length-md.offset, &elems)
	}
	md.body.RichText = &pb.IMMsgBody_RichText{
		Elems: elems,
	}
	return nil
}

func (md *messageDecoder) decodeText(length int64, elems *[]*pb.IMMsgBody_Elem) error {
	text := make([]byte, length)
	_, _ = md.buffer.Read(text)
	*elems = append(*elems, &pb.IMMsgBody_Elem{
		Text: &pb.IMMsgBody_Text{
			Str: text,
		},
	})
	md.offset += length
	return nil
}

func (md *messageDecoder) decodeFace(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem) error {
	_, _ = md.buffer.Read(make([]byte, entity.Length))
	url, _ := url.Parse(entity.URL)
	query := url.Query()
	id, _ := strconv.ParseInt(query.Get("id"), 10, 32)
	pid, _ := base64.RawStdEncoding.DecodeString(query.Get("pid"))
	sid, _ := base64.RawStdEncoding.DecodeString(query.Get("sid"))
	text := face.FaceType(id).String()
	if len(pid)+len(sid) > 0 {
		buf, _ := proto.Marshal(&pb.HummerCommonElement_MsgElemInfoServtype37{
			Packid:      pid,
			Stickerid:   sid,
			Qsid:        uint32(id),
			Sourcetype:  1,
			Stickertype: 1,
			Text:        []byte(text),
		})
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			CommonElem: &pb.IMMsgBody_CommonElem{
				ServiceType:  37,
				PbElem:       buf,
				BusinessType: 1,
			},
		})
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			Text: &pb.IMMsgBody_Text{
				Str: []byte("[" + strings.TrimLeft(text, "/") + "]请使用最新版手机QQ体验新功能"),
			},
		})
	} else if id < 260 {
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			Face: &pb.IMMsgBody_Face{
				Index: uint32(id),
			},
		})
	} else {
		buf, _ := proto.Marshal(&pb.HummerCommonElement_MsgElemInfoServtype33{
			Index:  uint32(id),
			Text:   []byte(text),
			Compat: []byte(text),
		})
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			CommonElem: &pb.IMMsgBody_CommonElem{
				ServiceType:  33,
				PbElem:       buf,
				BusinessType: 1,
			},
		})
	}
	md.offset += entity.Length
	return nil
}

func (md *messageDecoder) decodePhoto(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem, msg *penguin.Message) error {
	_, _ = md.buffer.Read(make([]byte, entity.Length))
	url, _ := url.Parse(entity.URL)
	query := url.Query()
	photo := msg.Photo
	md5, _ := hex.DecodeString(query.Get("md5"))
	if photo != nil && reflect.DeepEqual(md5, photo.Digest.MD5) {
		if msg.Chat.Type == penguin.ChatTypeGroup {
			*elems = append(*elems, &pb.IMMsgBody_Elem{
				CustomFace: &pb.IMMsgBody_CustomFace{
					FilePath: photo.Name,
					FileId:   uint32(photo.ID),
					FileType: ParsePhotoType(path.Ext(photo.Name)),
					Useful:   1,
					Md5:      photo.Digest.MD5,
					BizType:  0,
					Width:    uint32(photo.Width),
					Height:   uint32(photo.Height),
					Size:     uint32(photo.Size),
					Origin:   1,
				},
			})
		} else {
			*elems = append(*elems, &pb.IMMsgBody_Elem{
				Text: &pb.IMMsgBody_Text{
					Str: []byte(msg.Photo.Name),
				},
			})
		}
	} else {
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			Text: &pb.IMMsgBody_Text{
				Str: []byte("not match: " + msg.Photo.Name),
			},
		})
	}
	md.offset += entity.Length
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
	me.next = 0
	elems := me.body.GetRichText().GetElems()
	for i, elem := range elems {
		if v := elem.GetSrcMsg(); v != nil {
			me.encodeSrcMsg(v, msg)
			if i == 0 {
				me.next++
			}
			if me.next < len(elems) {
				if v := elems[me.next].GetText(); v != nil && len(v.GetAttr_6Buf()) != 0 {
					me.next++
				}
			}
			if me.next < len(elems) {
				if v := elems[me.next].GetText(); v != nil {
					if len(v.GetStr()) == 1 && v.GetStr()[0] == ' ' {
						me.next++
					} else if strings.HasPrefix(string(v.GetStr()), " ") {
						v.Str = v.GetStr()[1:]
					}
				}
			}
			break
		}
	}
	for ; me.next < len(elems); me.next++ {
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

func (me *messageEncoder) encodeNotOnlineImage(elem *pb.IMMsgBody_NotOnlineImage, msg *penguin.Message, flash ...bool) {
	n, _ := me.buffer.Write([]byte("[图片]"))
	entity := &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x&uid=%d", elem.GetPicMd5(), msg.From.Account.ID),
	}
	if len(flash) > 0 && flash[0] {
		entity.URL += "&flash=true"
	}
	msg.Entities = append(msg.Entities, entity)
	me.offset += int64(n)
}

func (me *messageEncoder) encodeCustomFace(elem *pb.IMMsgBody_CustomFace, msg *penguin.Message, flash ...bool) {
	n, _ := me.buffer.Write([]byte("[图片]"))
	entity := &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x&uid=%d", elem.GetMd5(), msg.From.Account.ID),
	}
	if len(flash) > 0 && flash[0] {
		entity.URL += "&flash=true"
	}
	msg.Entities = append(msg.Entities, entity)
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
	case 3: // photo
		info := pb.HummerCommonElement_MsgElemInfoServtype3{}
		_ = proto.Unmarshal(elem.GetPbElem(), &info)
		if v := info.GetFlashTroopPic(); v != nil {
			me.encodeCustomFace(v, msg, true)
		} else if v := info.GetFlashC2CPic(); v != nil {
			me.encodeNotOnlineImage(v, msg, true)
		}
		me.next++
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

func (me *messageEncoder) encodeSrcMsg(elem *pb.IMMsgBody_SourceMsg, msg *penguin.Message) {
	// TODO: implement
	msg.ReplyTo = &penguin.Message{
		MessageID: int64(elem.GetOrigSeqs()[0]),
		Chat:      msg.Chat,
		From:      &penguin.User{Account: &penguin.Account{ID: int64(elem.GetSenderUin())}},
		Time:      int64(elem.GetTime()),
		Text:      "",
		Entities:  []*penguin.MessageEntity{},
	}
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
