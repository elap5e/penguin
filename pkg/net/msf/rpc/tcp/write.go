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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (c *codec) WriteRequest(req *rpc.Request, args *rpc.Args) error {
	if req.Version != rpc.VersionDefault && req.Version != rpc.VersionSimple {
		return fmt.Errorf("tcp: unsupported version 0x%x", req.Version)
	}

	fake, session, tickets := c.cl.GetFakeSource(args.Uin), c.cl.GetSession(args.Uin), c.cl.GetTickets(args.Uin)
	p, _ := json.MarshalIndent(session, "", "  ")
	log.Trace("codec.write session:\n%s", string(p))
	p, _ = json.MarshalIndent(tickets, "", "  ")
	log.Trace("codec.write tickets:\n%s", string(p))
	body := bytes.NewBuffer([]byte{})
	body.WriteUint32(0)
	if req.Version == rpc.VersionDefault {
		body.WriteInt32(args.Seq)
		body.WriteInt32(fake.App.FixID)
		body.WriteInt32(fake.App.AppID)
		tmp := make([]byte, 12)
		tmp[0x0] = fake.Device.NetworkType
		tmp[0xa] = fake.Device.NetIPFamily
		body.Write(tmp)
		body.WriteBytesL32(tickets.A2.Sig)
	}
	body.WriteStringL32(args.ServiceMethod)
	body.WriteBytesL32(session.Cookie)
	if req.Version == rpc.VersionDefault {
		body.WriteStringL32(fake.Device.IMEI)
		body.WriteBytesL32(tickets.KSID)
		body.WriteStringL16("|" + fake.Device.IMSI + "|A" + fake.App.Version + "." + fake.App.Revision)
	}
	body.WriteBytesL32(args.ReserveField)
	body.WriteUint32At(uint32(body.Len()), 0)
	body.WriteBytesL32(args.Payload)
	log.Trace("codec.write.body:\n%s", hex.Dump(body.Bytes()))

	method := strings.ToLower(req.ServiceMethod)
	if method == "heartbeat.ping" || method == "heartbeat.alive" || method == "client.correcttime" {
		req.EncryptType = rpc.EncryptTypeNotNeedEncrypt
	} else {
		cipher, d2 := tea.NewCipher([16]byte{}), tickets.D2
		if len(d2.Sig) == 0 ||
			method == "login.auth" ||
			method == "login.chguin" ||
			method == "grayuinpro.check" ||
			method == "wtlogin.login" ||
			method == "wtlogin.name2uin" ||
			method == "wtlogin.exchange_emp" ||
			method == "wtlogin.trans_emp" ||
			method == "account.requestverifywtlogin_emp" ||
			method == "account.requestrebindmblwtlogin_emp" ||
			method == "connauthsvr.get_app_info_emp" ||
			method == "connauthsvr.get_auth_api_list_emp" ||
			method == "connauthsvr.sdk_auth_api_emp" ||
			method == "qqconnectlogin.pre_auth_emp" ||
			method == "qqconnectlogin.auth_emp" {
			req.EncryptType = rpc.EncryptTypeEncryptByZeros
		} else {
			cipher.SetKey(d2.Key.Get())
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
			head.WriteBytesL32(tickets.D2.Sig)
		} else {
			head.WriteUint32(4)
		}
	case rpc.VersionSimple:
		head.WriteInt32(args.Seq)
	}
	head.WriteByte(0)
	head.WriteStringL32(req.Username)
	head.WriteUint32At(uint32(head.Len()+body.Len()), 0)

	log.Trace("codec.write:\n%s", hex.Dump(append(head.Bytes(), body.Bytes()...)))
	if _, err := head.WriteTo(c.conn); err != nil {
		return err
	}
	if _, err := body.WriteTo(c.conn); err != nil {
		return err
	}
	return nil
}
