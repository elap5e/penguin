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

package rpc

import (
	"errors"
)

type Call struct {
	ServiceMethod string
	Seq           int32
	Version       uint32
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
	Version       uint32 `json:"version,omitempty"`
	Uin           int64  `json:"uin,omitempty"`
	Seq           int32  `json:"seq,omitempty"`
	ServiceMethod string `json:"service_method,omitempty"`
	ReserveField  []byte `json:"reserve_field,omitempty"`
	Payload       []byte `json:"-"`
}

type Reply struct {
	Version       uint32 `json:"version,omitempty"`
	Uin           int64  `json:"uin,omitempty"`
	Seq           int32  `json:"seq,omitempty"`
	Code          int32  `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	ServiceMethod string `json:"service_method,omitempty"`
	Cookie        []byte `json:"cookie,omitempty"`
	Flag          uint32 `json:"flag,omitempty"`
	ReserveField  []byte `json:"reserve_field,omitempty"`
	Payload       []byte `json:"-"`
}

type Request struct {
	ServiceMethod string `json:"service_method,omitempty"`
	Seq           int32  `json:"seq,omitempty"`
	Version       uint32 `json:"version,omitempty"`
	EncryptType   uint8  `json:"encrypt_type,omitempty"`
	Username      string `json:"username,omitempty"`
}

type Response struct {
	ServiceMethod string `json:"service_method,omitempty"`
	Seq           int32  `json:"seq,omitempty"`
	Version       uint32 `json:"version,omitempty"`
	EncryptType   uint8  `json:"encrypt_type,omitempty"`
	Username      string `json:"username,omitempty"`
}

var (
	ErrCachedPush = errors.New("cached push")
	ErrNotHandled = errors.New("not handled")

	ErrHeartbeatTimeout = errors.New("heartbeat timeout")

	ErrWriteClosed  = errors.New("write closed")
	ErrWriteTimeout = errors.New("write timeout")
)

var (
	ContextErrorKey = struct{}{}
)
