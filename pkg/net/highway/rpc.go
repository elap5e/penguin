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

package highway

import (
	"github.com/elap5e/penguin/pkg/net/highway/pb"
)

type Call struct {
	ServiceMethod string
	Seq           int32
	Args          *Args
	Reply         *Reply
	Error         error
	Done          chan *Call
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
	}
}

type Args struct {
	Uin           int64                         `json:"uin,omitempty"`
	Seq           int32                         `json:"seq,omitempty"`
	ServiceMethod string                        `json:"service_method,omitempty"`
	CommandID     int32                         `json:"command_id,omitempty"`
	SegHead       *pb.CSDataHighwayHead_SegHead `json:"seg_head,omitempty"`
	Payload       []byte                        `json:"-"`
}

type Reply struct {
	Uin           int64  `json:"uin,omitempty"`
	Seq           int32  `json:"seq,omitempty"`
	Code          int32  `json:"code,omitempty"`
	ServiceMethod string `json:"service_method,omitempty"`
	Payload       []byte `json:"-"`
}

type Request struct {
	Username      string `json:"username,omitempty"`
	Seq           int32  `json:"seq,omitempty"`
	ServiceMethod string `json:"service_method,omitempty"`
	CommandID     int32  `json:"command_id,omitempty"`
	AppID         int32  `json:"app_id,omitempty"`
}

type Response struct {
	Username      string `json:"username,omitempty"`
	Seq           int32  `json:"seq,omitempty"`
	ServiceMethod string `json:"service_method,omitempty"`
	Code          int32  `json:"code,omitempty"`
}

const (
	ServiceMethodEcho   = "PicUp.Echo"
	ServiceMethodUpload = "PicUp.DataUp"
)
