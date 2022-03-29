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

package message

import (
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/daemon/service"
)

type Daemon interface {
	GetServiceManager() *service.Manager
	OnRecvMessage(head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) error
}

type Message struct {
	FromUin         int64            `jce:"0" json:"from_uin,omitempty"`
	Time            int64            `jce:"1" json:"time,omitempty"`
	Type            int16            `jce:"2" json:"type,omitempty"`
	Seq             int32            `jce:"3" json:"seq,omitempty"`
	Message         string           `jce:"4" json:"message,omitempty"`
	RealMessageTime int64            `jce:"5" json:"real_message_time,omitempty"`
	MessageBytes    []byte           `jce:"6" json:"message_bytes,omitempty"`
	AppShareID      int64            `jce:"7" json:"app_share_id,omitempty"`
	MessageCookie   []byte           `jce:"8" json:"message_cookie,omitempty"`
	AppShareCookie  []byte           `jce:"9" json:"app_share_cookie,omitempty"`
	MessageUid      int64            `jce:"10" json:"message_uid,omitempty"`
	LastChangeTime  int64            `jce:"11" json:"last_change_time,omitempty"`
	Pictures        []*Picture       `jce:"12" json:"pictures,omitempty"`
	ShareData       *ShareData       `jce:"13" json:"share_data,omitempty"`
	FromInstanceID  int64            `jce:"14" json:"from_instance_id,omitempty"`
	FromRemark      []byte           `jce:"15" json:"from_remark,omitempty"`
	FromMobile      string           `jce:"16" json:"from_mobile,omitempty"`
	FromName        string           `jce:"17" json:"from_name,omitempty"`
	FromNick        []string         `jce:"18" json:"from_nick,omitempty"`
	TempMessageHead *TempMessageHead `jce:"19" json:"temp_message_head,omitempty"`
}

type MessageDelete struct {
	FromUin       int64  `jce:"0" json:"from_uin,omitempty"`
	MessageTime   int64  `jce:"1" json:"message_time,omitempty"`
	MessageSeq    int32  `jce:"2" json:"message_seq,omitempty"`
	MessageCookie []byte `jce:"3" json:"message_cookie,omitempty"`
}

type MessageReadedC2C struct {
	Uin          int64 `jce:"0" json:"uin,omitempty"`
	LastReadTime int64 `jce:"1" json:"last_read_time,omitempty"`
}

type MessageReadedGroup struct {
	Uin        int64 `jce:"0" json:"uin,omitempty"`
	Type       int64 `jce:"1" json:"type,omitempty"`
	MemberSeq  int32 `jce:"2" json:"member_seq,omitempty"`
	MessageSeq int32 `jce:"3" json:"message_seq,omitempty"`
}

type MessageReadedDiscuss struct {
	Uin        int64 `jce:"0" json:"uin,omitempty"`
	Type       int64 `jce:"1" json:"type,omitempty"`
	MemberSeq  int32 `jce:"2" json:"member_seq,omitempty"`
	MessageSeq int32 `jce:"3" json:"message_seq,omitempty"`
}

type Picture struct {
	Path []byte `jce:"0" json:"path,omitempty"`
	Host []byte `jce:"1" json:"host,omitempty"`
}

type ShareData struct {
	PackageName string `jce:"0" json:"package_name,omitempty"`
	MessageTail string `jce:"1" json:"message_tail,omitempty"`
	PictureURL  string `jce:"2" json:"picture_url,omitempty"`
	URL         string `jce:"3" json:"url,omitempty"`
}

type TempMessageHead struct {
	C2CType     int32 `jce:"0" json:"c2c_type,omitempty"`
	ServiceType int32 `jce:"1" json:"service_type,omitempty"`
}

type UinPairMessage struct {
	LastReadTime     int64      `jce:"1" json:"last_read_time,omitempty"`
	PeerUin          int64      `jce:"2" json:"peer_uin,omitempty"`
	MessageCompleted int64      `jce:"3" json:"message_completed,omitempty"`
	Messages         []*Message `jce:"4" json:"messages,omitempty"`
}
