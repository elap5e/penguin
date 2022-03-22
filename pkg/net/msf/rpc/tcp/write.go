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

package tcp

import (
	"fmt"
	"strings"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (c *codec) WriteRequest(req *rpc.Request, args *rpc.Args) error {
	if req.Version != rpc.VersionDefault && req.Version != rpc.VersionSimple {
		return fmt.Errorf("tcp: unsupported version 0x%x", req.Version)
	}

	body := bytes.NewBuffer([]byte{})
	body.WriteUint32(0)
	if req.Version == rpc.VersionDefault {
		body.WriteInt32(args.Seq)
		body.WriteInt32(args.FixID)
		body.WriteInt32(args.AppID)
		tmp := make([]byte, 12)
		tmp[0x0] = c.c.GetFakeApp(req.Username).NetworkType
		tmp[0xa] = c.c.GetFakeApp(req.Username).NetIPFamily
		body.Write(tmp)
		body.WriteBytesL32(c.c.GetTickets(req.Username).A2().Sig())
	}
	body.WriteStringL32(args.ServiceMethod)
	body.WriteBytesL32(args.Cookie)
	if req.Version == rpc.VersionDefault {
		body.WriteStringL32(c.c.GetFakeApp(req.Username).IMEI)
		body.WriteBytesL32(c.c.GetFakeApp(req.Username).KSID)
		body.WriteStringL16("" + "|" + c.c.GetFakeApp(req.Username).IMSI + "|A" + c.c.GetFakeApp(req.Username).Revision)
	}
	body.WriteBytesL32(args.ReserveField)
	body.WriteUint32At(uint32(body.Len()), 0)
	body.WriteBytesL32(args.Payload)

	method := strings.ToLower(req.ServiceMethod)
	if method == "heartbeat.ping" || method == "heartbeat.alive" || method == "client.correcttime" {
		req.EncryptType = rpc.EncryptTypeNotNeedEncrypt
	} else {
		cipher, d2 := tea.NewCipher([16]byte{}), c.c.GetTickets(req.Username).D2()
		if len(d2.Sig()) == 0 ||
			method == "login.auth" ||
			method == "login.chguin" ||
			method == "grayuinpro.check" ||
			method == "wtlogin.login" ||
			method == "wtlogin.name2uin" ||
			method == "wtlogin.exchange_emp" ||
			method == "wtlogin.trans_emp" ||
			method == "account.requestverifywtlogin_emp" ||
			method == "account.requestrebindmblwtLogin_emp" ||
			method == "connauthsvr.get_app_info_emp" ||
			method == "connauthsvr.get_auth_api_list_emp" ||
			method == "connauthsvr.sdk_auth_api_emp" ||
			method == "qqconnectlogin.pre_auth_emp" ||
			method == "qqconnectlogin.auth_emp" {
			req.EncryptType = rpc.EncryptTypeEncryptByZeros
		} else {
			cipher.SetKey(d2.Key())
			req.EncryptType = rpc.EncryptTypeEncryptByD2Key
		}
		body = bytes.NewBuffer(cipher.Encrypt(body.Bytes()))
	}

	head := bytes.NewBuffer([]byte{})
	head.Reset()
	head.WriteUint32(0)
	head.WriteUint32(req.Version)
	head.WriteByte(req.EncryptType)
	switch req.Version {
	case rpc.VersionDefault:
		if req.EncryptType == rpc.EncryptTypeEncryptByD2Key {
			head.WriteBytesL32(c.c.GetTickets(req.Username).D2().Sig())
		} else {
			head.WriteUint32(4)
		}
	case rpc.VersionSimple:
		head.WriteInt32(args.Seq)
	}
	head.WriteByte(0)
	head.WriteStringL32(req.Username)
	head.WriteUint32At(uint32(head.Len()+body.Len()), 0)

	if _, err := head.WriteTo(c.conn); err != nil {
		return err
	}
	if _, err := body.WriteTo(c.conn); err != nil {
		return err
	}
	return nil
}
