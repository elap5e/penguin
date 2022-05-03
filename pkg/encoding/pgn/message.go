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
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/message/face"
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/log"
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
		if entity.Type == "mention" {
			md.decodeMention(entity, &elems)
		} else if entity.Type == "face" {
			md.decodeFace(entity, &elems)
		} else if entity.Type == "dice" {
			md.decodeDice(entity, &elems)
		} else if entity.Type == "photo" {
			md.decodePhoto(entity, &elems, msg)
		} else if entity.Type == "video" {
			md.decodeVideo(entity, &elems, msg)
		} else if entity.Type == "voice" {
			md.decodeVoice(entity, &elems, msg)
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
	// a text entity is no more than 4500 runes
	*elems = append(*elems, &pb.IMMsgBody_Elem{
		Text: &pb.IMMsgBody_Text{
			Str: text,
		},
	})
	md.offset += length
	return nil
}

func (md *messageDecoder) decodeMention(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem) error {
	text := make([]byte, entity.Length)
	_, _ = md.buffer.Read(text)
	url, _ := url.Parse(entity.URL)
	query := url.Query()
	buf := make([]byte, 13)
	buf[1] = 0x01
	binary.BigEndian.PutUint16(buf[4:], uint16(len([]rune(string(text)))))
	id, _ := strconv.ParseUint(query.Get("id"), 10, 64)
	if id == 0 {
		buf[6] = 0x01
	} else {
		binary.BigEndian.PutUint32(buf[7:], uint32(id))
	}
	*elems = append(*elems, &pb.IMMsgBody_Elem{
		Text: &pb.IMMsgBody_Text{
			Str:       text,
			Attr_6Buf: buf,
		},
	})
	md.offset += entity.Length
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
		info := []byte("[" + strings.TrimLeft(text, "/") + "]请使用最新版手机QQ体验新功能")
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			Text: &pb.IMMsgBody_Text{
				Str:       []byte(text),
				PbReserve: append([]byte{0x0a, byte(len(info))}, info...),
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

func (md *messageDecoder) decodeDice(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem) error {
	text := make([]byte, entity.Length)
	_, _ = md.buffer.Read(text)
	url, _ := url.Parse(entity.URL)
	query := url.Query()
	key := query.Get("key")
	value, _ := strconv.Atoi(query.Get("value"))
	if key == "409e2a69b16918f9" { // dice
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			MarketFace: &pb.IMMsgBody_MarketFace{
				FaceName:    text,
				ItemType:    6,
				FaceInfo:    1,
				FaceId:      []byte{0x48, 0x23, 0xd3, 0xad, 0xb1, 0x5d, 0xf0, 0x80, 0x14, 0xce, 0x5d, 0x67, 0x96, 0xb7, 0x6e, 0xe1},
				TabId:       11464,
				SubType:     3,
				Key:         []byte(key),
				ImageWidth:  200,
				ImageHeight: 200,
				Mobileparam: []byte("rscType?1;value=" + strconv.Itoa(value)),
				PbReserve:   []byte{0x0a, 0x06, 0x08, 0xc8, 0x01, 0x10, 0xc8, 0x01, 0x40, 0x01},
			},
		})
	} else if key == "7de39febcf45e6db" { // rock-paper-scissors
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			MarketFace: &pb.IMMsgBody_MarketFace{
				FaceName:    text,
				ItemType:    6,
				FaceInfo:    1,
				FaceId:      []byte{0x83, 0xc8, 0xa2, 0x93, 0xae, 0x65, 0xca, 0x14, 0x0f, 0x34, 0x81, 0x20, 0xa7, 0x74, 0x48, 0xee},
				TabId:       11415,
				SubType:     3,
				Key:         []byte(key),
				ImageWidth:  200,
				ImageHeight: 200,
				Mobileparam: []byte("rscType?1;value=" + strconv.Itoa(value)),
				PbReserve:   []byte{0x0a, 0x06, 0x08, 0xc8, 0x01, 0x10, 0xc8, 0x01, 0x40, 0x01},
			},
		})
	}
	*elems = append(*elems, &pb.IMMsgBody_Elem{
		Text: &pb.IMMsgBody_Text{
			Str: text,
		},
	})
	md.offset += entity.Length
	return nil
}

func (md *messageDecoder) decodePhoto(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem, msg *penguin.Message) error {
	_, _ = md.buffer.Read(make([]byte, entity.Length))
	url, _ := url.Parse(entity.URL)
	query := url.Query()
	photo := msg.Photo
	md5, _ := hex.DecodeString(query.Get("md5"))
	if photo != nil && !reflect.DeepEqual(md5, photo.Digest.MD5) {
		*elems = append(*elems, &pb.IMMsgBody_Elem{
			Text: &pb.IMMsgBody_Text{
				Str: []byte("not match: " + msg.Photo.Name),
			},
		})
	} else {
		if photo == nil {
			w, _ := strconv.Atoi(query.Get("w"))
			h, _ := strconv.Atoi(query.Get("h"))
			size, _ := strconv.ParseInt(query.Get("size"), 10, 64)
			if w == 0 {
				w = 640
			}
			if h == 0 {
				h = 480
			}
			if size == 0 {
				size = int64(w * h)
			}
			photo = &penguin.Photo{
				ID:     random.Int63n(0xffffffff),
				Name:   hex.EncodeToString(md5) + ".jpg",
				Size:   size,
				Width:  w,
				Height: h,
				Digest: &penguin.Digest{MD5: md5},
			}
		}
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
					Str: []byte(hex.EncodeToString(md5)),
				},
			})
		}
	}
	md.offset += entity.Length
	return nil
}

func (md *messageDecoder) decodeLightApp(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem, msg *penguin.Message) error {
	return nil
}

func (md *messageDecoder) decodeVideo(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem, msg *penguin.Message) error {
	*elems = append(*elems, &pb.IMMsgBody_Elem{
		Text: &pb.IMMsgBody_Text{
			Str: []byte("你的QQ暂不支持查看视频短片，请期待后续版本。"),
		},
	})
	return nil
}

func (md *messageDecoder) decodeVoice(entity *penguin.MessageEntity, elems *[]*pb.IMMsgBody_Elem, msg *penguin.Message) error {
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
	defer func() {
		msg.Text = me.buffer.String()
	}()

	if ptt := me.body.GetRichText().GetPtt(); ptt != nil {
		me.encodeVoice(ptt, msg)
		return nil
	}

	elems := me.body.GetRichText().GetElems()
	for i, elem := range elems {
		if v := elem.GetVideoFile(); v != nil {
			me.encodeVideo(v, msg)
			return nil
		} else if v := elem.GetLightApp(); v != nil {
			me.encodeLightApp(v, msg)
			return nil
		} else if v := elem.GetRichMsg(); v != nil {
			me.encodeRichMsg(v, msg)
			return nil
		} else if v := elem.GetQqwalletMsg(); v != nil {
			me.encodeWallet(v, msg)
			return nil
		} else if v := elem.GetSrcMsg(); v != nil {
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
		if v := elem.GetAnonGroupMsg(); v != nil {
			me.encodeAnonGroupMsg(v, msg)
		} else if v := elem.GetShakeWindow(); v != nil {
			me.encodeShakeWindow(v, msg)
		} else if v := elem.GetText(); v != nil {
			me.encodeText(v, msg)
		} else if v := elem.GetFace(); v != nil {
			me.encodeFace(v, msg)
		} else if v := elem.GetMarketFace(); v != nil {
			me.encodeMarketFace(v, msg)
		} else if v := elem.GetCustomFace(); v != nil {
			me.encodeCustomFace(v, msg)
		} else if v := elem.GetNotOnlineImage(); v != nil {
			me.encodeNotOnlineImage(v, msg)
		} else if v := elem.GetCommonElem(); v != nil {
			me.encodeCommonElem(v, msg)
		}
	}
	return nil
}

func (me *messageEncoder) encodeVoice(elem *pb.IMMsgBody_Ptt, msg *penguin.Message) {
	msg.Voice = &penguin.Voice{
		ID:   int64(elem.GetFileId()),
		Path: string(elem.GetDownPara()),
		Name: hex.EncodeToString(elem.GetFileMd5()) + ".amr",
		Size: int64(elem.GetFileSize()),
		Digest: &penguin.Digest{
			MD5: elem.GetFileMd5(),
		},
	}
	n, _ := me.buffer.Write([]byte("[语音]"))
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "voice",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x", elem.GetFileMd5()),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeVideo(elem *pb.IMMsgBody_VideoFile, msg *penguin.Message) {
	p, err := hex.DecodeString(string(elem.GetFileUuid()))
	if err != nil {
		log.Warn("hex.DecodeString(%s) error(%v)", elem.GetFileUuid(), err)
		p = append([]byte("~"), elem.GetFileUuid()...)
	}
	uuid := base64.RawURLEncoding.EncodeToString(p)
	msg.Video = &penguin.Video{
		UUID: uuid,
		Name: string(elem.GetFileName()),
		Size: int64(elem.GetFileSize()),
		Digest: &penguin.Digest{
			MD5: elem.GetFileMd5(),
		},
	}
	msg.Video.Thumbnail = &penguin.Photo{
		Size:   int64(elem.GetThumbFileSize()),
		Width:  int(elem.GetThumbWidth()),
		Height: int(elem.GetThumbHeight()),
		Digest: &penguin.Digest{
			MD5: elem.GetThumbFileMd5(),
		},
	}
	n, _ := me.buffer.Write([]byte("[视频]"))
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "video",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?md5=%x&uuid=%s", elem.GetFileMd5(), uuid),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeLightApp(elem *pb.IMMsgBody_LightAppElem, msg *penguin.Message) {
	data := elem.GetData()[1:]
	if elem.GetData()[0] == 1 {
		r, _ := zlib.NewReader(bytes.NewReader(data))
		defer r.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		data = buf.Bytes()
	}
	_, _ = me.buffer.Write(data)
}

func (me *messageEncoder) encodeRichMsg(elem *pb.IMMsgBody_RichMsg, msg *penguin.Message) {
	data := elem.GetTemplate_1()[1:]
	if elem.GetTemplate_1()[0] == 1 {
		r, _ := zlib.NewReader(bytes.NewReader(data))
		defer r.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		data = buf.Bytes()
	}
	_, _ = me.buffer.Write(data)
}

func (me *messageEncoder) encodeWallet(elem *pb.IMMsgBody_QQWalletMsg, msg *penguin.Message) {
	var n int
	if elem.GetAioBody().GetMsgtype() == 1 {
		n, _ = me.buffer.Write([]byte("[转账]"))
	} else if elem.GetAioBody().GetMsgtype() == 2 {
		n, _ = me.buffer.Write([]byte("[红包]"))
	} else {
		log.Warn("unknown wallet msgtype: %d", elem.GetAioBody().GetMsgtype())
	}
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "wallet",
		Offset: me.offset,
		Length: int64(n),
		URL:    "",
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeSrcMsg(elem *pb.IMMsgBody_SourceMsg, msg *penguin.Message) {
	msg.ReplyTo = &penguin.Message{
		MessageID: int64(elem.GetOrigSeqs()[0]),
		Chat:      msg.Chat,
		From:      &penguin.User{Account: &penguin.Account{ID: int64(elem.GetSenderUin())}},
		Time:      int64(elem.GetTime()),
	}
	_ = NewMessageEncoder(&pb.IMMsgBody_MsgBody{
		RichText: &pb.IMMsgBody_RichText{
			Elems: elem.GetElems(),
		},
	}).Encode(msg.ReplyTo)
	n, _ := me.buffer.Write([]byte("[回复:" + msg.ReplyTo.Text + "]:"))
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "reply",
		Offset: me.offset,
		Length: int64(n),
		URL:    fmt.Sprintf("?id=%d", msg.ReplyTo.MessageID),
	})
	me.offset += int64(n)
}

func (me *messageEncoder) encodeAnonGroupMsg(elem *pb.IMMsgBody_AnonymousGroupMsg, msg *penguin.Message) {
	msg.From.Account.Type = penguin.AccountTypeAnonymous
	msg.From.Display = string(elem.GetAnonNick())
	msg.Entities = append([]*penguin.MessageEntity{{
		Type:   "anonymous",
		Offset: me.offset,
		Length: 0,
		URL:    fmt.Sprintf("?id=%s", base64.RawURLEncoding.EncodeToString(elem.GetAnonId())),
	}}, msg.Entities...)
}

func (me *messageEncoder) encodeShakeWindow(elem *pb.IMMsgBody_ShakeWindow, msg *penguin.Message) {
	msg.Entities = append(msg.Entities, &penguin.MessageEntity{
		Type:   "poke",
		Offset: 0,
		Length: 0,
		URL:    "?id=0",
	})
}

func (me *messageEncoder) encodeText(elem *pb.IMMsgBody_Text, msg *penguin.Message) {
	n, _ := me.buffer.Write(elem.GetStr())
	if buf := elem.GetAttr_6Buf(); len(buf) > 12 {
		id := int64(buf[7])<<24 + int64(buf[8])<<16 + int64(buf[9])<<8 + int64(buf[10])
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
			Type:   "mention",
			Offset: me.offset,
			Length: int64(n),
			URL:    fmt.Sprintf("?id=%d", id),
			User:   &penguin.User{Account: &penguin.Account{ID: id}},
		})
	}
	if buf := elem.GetPbReserve(); len(buf) > 0 {
		info := pb.IMMsgBodyElem_TextPb{}
		_ = proto.Unmarshal(elem.GetPbReserve(), &info)
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
			Type:   "mention",
			Offset: me.offset,
			Length: int64(n),
			URL:    fmt.Sprintf("?id=%d", info.TinyId),
			User:   &penguin.User{Account: &penguin.Account{ID: info.TinyId}},
		})
	}
	me.offset += int64(n)
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

func (me *messageEncoder) encodeMarketFace(elem *pb.IMMsgBody_MarketFace, msg *penguin.Message) {
	me.next++
	name := elem.GetFaceName()
	if len(name) == 0 {
		name = me.body.GetRichText().GetElems()[me.next].GetText().GetStr()
	}
	n, _ := me.buffer.Write(name)
	key := string(elem.GetKey())
	if key == "409e2a69b16918f9" || // dice
		key == "7de39febcf45e6db" { // rock-paper-scissors
		value, _ := strconv.Atoi(strings.TrimPrefix(string(elem.GetMobileparam()), "rscType?1;value="))
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
			Type:   "dice",
			Offset: me.offset,
			Length: int64(n),
			URL:    fmt.Sprintf("?key=%s&value=%d", key, value),
		})
	} else {
		msg.Entities = append(msg.Entities, &penguin.MessageEntity{
			Type:   "sticker",
			Offset: me.offset,
			Length: int64(n),
			URL: fmt.Sprintf(
				"?id=%s&tid=%d&key=%s&h=%d&w=%d",
				base64.RawURLEncoding.EncodeToString(elem.GetFaceId()),
				elem.GetTabId(),
				key,
				elem.GetImageHeight(),
				elem.GetImageWidth(),
			),
		})
	}
	me.offset += int64(n)
}

func (me *messageEncoder) encodeCustomFace(elem *pb.IMMsgBody_CustomFace, msg *penguin.Message, flash ...bool) {
	msg.Photo = &penguin.Photo{
		ID:     int64(elem.GetFileId()),
		Name:   string(elem.GetGuid()),
		Size:   int64(elem.GetSize()),
		Width:  int(elem.GetWidth()),
		Height: int(elem.GetHeight()),
		Digest: &penguin.Digest{
			MD5: elem.GetMd5(),
		},
	}
	n, _ := me.buffer.Write([]byte("[图片]"))
	entity := &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL: fmt.Sprintf(
			"?md5=%x&size=%d&w=%d&h=%d&uid=%d",
			elem.GetMd5(),
			elem.GetSize(),
			elem.GetWidth(),
			elem.GetHeight(),
			msg.From.Account.ID,
		),
	}
	if len(flash) > 0 && flash[0] {
		entity.URL += "&flash=true"
	}
	msg.Entities = append(msg.Entities, entity)
	me.offset += int64(n)
}

func (me *messageEncoder) encodeNotOnlineImage(elem *pb.IMMsgBody_NotOnlineImage, msg *penguin.Message, flash ...bool) {
	msg.Photo = &penguin.Photo{
		ID:     int64(elem.GetFileId()),
		Size:   int64(elem.GetFileLen()),
		Width:  int(elem.GetPicWidth()),
		Height: int(elem.GetPicHeight()),
		Digest: &penguin.Digest{
			MD5: elem.GetPicMd5(),
		},
	}
	n, _ := me.buffer.Write([]byte("[图片]"))
	entity := &penguin.MessageEntity{
		Type:   "photo",
		Offset: me.offset,
		Length: int64(n),
		URL: fmt.Sprintf(
			"?md5=%x&size=%d&w=%d&h=%d&uid=%d",
			elem.GetPicMd5(),
			elem.GetFileLen(),
			elem.GetPicWidth(),
			elem.GetPicHeight(),
			msg.From.Account.ID,
		),
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
