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

package auth

import (
	"context"
	"crypto/rand"
	"strconv"

	"github.com/elap5e/penguin/daemon/auth/pb"
	"github.com/elap5e/penguin/pkg/encoding/oicq"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Manager struct {
	ctx context.Context

	c rpc.Client

	// session
	mapExtraData map[int64]*ExtraData
}

func NewManager(ctx context.Context, c rpc.Client) *Manager {
	return &Manager{
		ctx:          ctx,
		c:            c,
		mapExtraData: make(map[int64]*ExtraData),
	}
}

type CaptchaSign struct {
	Ticket string `json:"ticket"`
	Random string `json:"random"`
	Return string `json:"return"`
	AppID  uint64 `json:"app_id"`
}

type ExtraData struct {
	Login []byte `json:"login,omitempty"`
	T16A  []byte `json:"t16a,omitempty"`
	T172  []byte `json:"t172,omitempty"`

	SessionAuth  []byte `json:"session_auth,omitempty"`
	PictureSign  []byte `json:"picture_sign,omitempty"`
	PictureData  []byte `json:"picture_data,omitempty"`
	CaptchaSign  string `json:"captcha_sign,omitempty"`
	ErrorCode    uint32 `json:"error_code,omitempty"`
	ErrorTitle   string `json:"error_title,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Message      string `json:"message,omitempty"`
	SMSMobile    string `json:"sms_mobile,omitempty"`

	Salt uint64 `json:"salt,omitempty"`

	SignInCodeSign []byte                     `json:"signin_code_sign,omitempty"`
	ThirdPartLogin *pb.ThirdPartLogin_RspBody `json:"third_part_login,omitempty"`

	T119 []byte          `json:"t119,omitempty"`
	T150 []byte          `json:"t150,omitempty"`
	T161 []byte          `json:"t161,omitempty"`
	T174 []byte          `json:"t174,omitempty"`
	T17B []byte          `json:"t17b,omitempty"`
	T182 []byte          `json:"t182,omitempty"`
	T401 *rpc.Key16Bytes `json:"t401,omitempty"`
	T402 []byte          `json:"t402,omitempty"`
	T403 []byte          `json:"t403,omitempty"`
	T542 []byte          `json:"t542,omitempty"`
	T546 []byte          `json:"t546,omitempty"`
	T547 []byte          `json:"t547,omitempty"`
}

type Request struct {
	ServiceMethod string     `json:"service_method,omitempty"`
	Seq           int32      `json:"seq,omitempty"`
	Data          *oicq.Data `json:"data,omitempty"`
}

type Response struct {
	ServiceMethod string     `json:"service_method,omitempty"`
	Seq           int32      `json:"seq,omitempty"`
	Data          *oicq.Data `json:"data,omitempty"`
	ExtraData     *ExtraData `json:"extra_data,omitempty"`
}

func checkUsername(username string) bool {
	uin, err := strconv.Atoi(username)
	if err != nil || uin < 10000 || uin > 4000000000 {
		return false
	}
	return true
}

func randomPassword() string {
	password := [16]byte{}
	rand.Read(password[:])
	strs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range password {
		password[i] = strs[password[i]%52]
	}
	return string(password[:])
}
