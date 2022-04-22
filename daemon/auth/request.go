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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
	"github.com/elap5e/penguin/pkg/encoding/oicq"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (m *Manager) request(req *Request) (*Response, error) {
	if req.Seq == 0 {
		req.Seq = m.c.GetNextSeq()
	}

	sess := m.c.GetSession(req.Data.Uin)
	req.Data.RandomKey = sess.RandomKey
	req.Data.KeyVersion = sess.KeyVersion
	req.Data.PublicKey = sess.PrivateKey.Public().Bytes()
	req.Data.SharedSecret = sess.SharedSecret

	resp, err := m.call(req)
	if err != nil {
		return nil, err
	}

	data, extraData := resp.Data, m.GetExtraData(resp.Data.Uin)
	switch data.Code {
	case 0x00:
		// VerifySignInCode
		if data.Type == 0x0012 {
			m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

			login := resp.ExtraData.ThirdPartLogin.MsgRspCmd_18.MsgRspPhoneSmsExtendLogin
			for k, v := range login.BindUinInfo {
				log.Printf("%d: %s(%s) photo:%s", k, v.Nick, v.MaskUin, v.ImageUrl)
			}
			log.Warn("auth unbinding notify: %s", login.UnbindWording)

			// select account and input password
			var code string
			fmt.Printf(">>> ")
			fmt.Scanln(&code)
			line, _ := strconv.Atoi(code)
			info := login.BindUinInfo[line]
			extraData.T542 = info.EncryptUin

			return m.name2Uin(info.Nick)
		}

		// success
		extraData.Login = resp.ExtraData.Login

		// decode t119
		tickets := m.c.GetTickets(data.Uin)
		key := [16]byte{}
		switch data.Type {
		default:
			log.Debug("auth unknown type: 0x%02x", data.Type)
			copy(key[:], tickets.A1.Key[:])
		case 0x0007: // VerifySMSCode
			copy(key[:], tickets.A1.Key[:])
		case 0x0009: // signInWithPassword
			copy(key[:], tickets.A1.Key[:])
		case 0x000a: // signInWithoutPassword.A2
			copy(key[:], tickets.A2.Key[:])
		case 0x000b: // signInWithoutPassword.D2
			key = md5.Sum(tickets.D2.Key[:])
		case 0x0014: // UnlockDevice
			copy(key[:], tickets.A1.Key[:])
		}
		t119, err := tea.NewCipher(key).Decrypt(resp.ExtraData.T119)
		if err != nil {
			return nil, err
		}

		tlvs := map[uint16]tlv.Codec{}
		buf := bytes.NewBuffer(t119)
		l, _ := buf.ReadUint16()
		for i := 0; i < int(l); i++ {
			v := tlv.TLV{}
			v.ReadFrom(buf)
			tlvs[v.GetType()] = &v
		}

		m.c.SetTickets(data.Uin, tlvs)
		m.c.SetSessionAuth(data.Uin, nil)
		return resp, nil
	case 0x02:
		// verify captcha or picture
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		var code string
		extraData.T546 = resp.ExtraData.T546
		extraData.T547 = resp.ExtraData.T547
		if resp.ExtraData.CaptchaSign != "" {
			l, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return nil, err
			}
			addr := l.Addr().(*net.TCPAddr).String()
			log.Warn("auth verify captcha, url: %s", strings.ReplaceAll(
				resp.ExtraData.CaptchaSign,
				"https://ssl.captcha.qq.com/template/wireless_mqq_captcha.html",
				"http://"+addr+"/index.html",
			))
			sign, err := serveHTTPVerifyCaptcha(l)
			if err != nil {
				return nil, err
			}
			return m.VerifyCaptcha(data.Uin, []byte(sign.Ticket))
		} else {
			log.Warn("auth verify picture")
			fmt.Printf(">>> ")
			fmt.Scanln(&code)
			return m.VerifyCaptcha(data.Uin, []byte(code), resp.ExtraData.PictureSign)
		}
	case 0xa0:
		// verify sms code
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.T17B = resp.ExtraData.T17B
		log.Warn("auth verify sms code")
		var code string
		fmt.Printf(">>> ")
		fmt.Scanln(&code)
		return m.VerifySMSCode(data.Uin, []byte(code))
	case 0xd0:
		// verify signin code
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.SignInCodeSign = resp.ExtraData.SignInCodeSign
		extraData.T182 = resp.ExtraData.T182
		extraData.Salt = resp.ExtraData.Salt
		log.Warn("auth verify signin code")
		var code string
		fmt.Printf(">>> ")
		fmt.Scanln(&code)
		return m.VerifySignInCode([]byte(code))
	case 0xcc:
		// unlock device
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.T402 = resp.ExtraData.T402
		extraData.T403 = resp.ExtraData.T403
		extraData.T401.Set(md5.Sum(
			append(append(
				m.c.GetFakeSource(data.Uin).Device.GUID[:],
				sess.RandomPass[:]...),
				extraData.T402...),
		))
		log.Warn("auth unlock device")
		return m.unlockDevice(data.Uin)
	case 0xef:
		// resend sms code
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.T174 = resp.ExtraData.T174
		extraData.T402 = resp.ExtraData.T402
		extraData.T403 = resp.ExtraData.T403
		extraData.T401.Set(md5.Sum(
			append(append(
				m.c.GetFakeSource(data.Uin).Device.GUID[:],
				sess.RandomPass[:]...),
				extraData.T402...),
		))
		log.Warn("auth resend sms code, mobile:", resp.ExtraData.SMSMobile)
		return m.resendSMSCode(data.Uin)
	case 0x01:
		log.Debug("invalid login, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("invalid password")
	case 0x06:
		log.Debug("not implement, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x09:
		log.Debug("invalid service, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x0a:
		log.Debug("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("service temporarily unavailable")
	case 0x10:
		log.Debug("session expired, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x28:
		log.Debug("protection mode, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x9a:
		log.Debug("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("service temporarily unavailable")
	case 0xa1:
		log.Debug("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("too many sms verify requests")
	case 0xa2:
		log.Debug("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("frequent sms verify requests")
	case 0xa4:
		log.Debug("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("bad requests")
	case 0xd5:
		log.Debug("phone number not valid, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("phone number not valid")
	case 0xdb:
		log.Debug("phone number not registered, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("phone number not registered")
	case 0xeb:
		log.Debug("version too low, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("too many failures")
	case 0xed:
		log.Debug("invalid device, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("too many failures")
	default:
		log.Debug("unknown code: 0x%02x", data.Code)
		return nil, fmt.Errorf("unknown code: 0x%02x", data.Code)
	}
	return resp, nil
}

func (m *Manager) call(req *Request) (*Response, error) {
	p, _ := json.Marshal(req)
	log.Trace("uin:%d seq:%d cmd:%s meta:%s", req.Data.Uin, req.Seq, req.ServiceMethod, string(p))
	p, err := oicq.Marshal(req.Data)
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{Uin: req.Data.Uin, Seq: req.Seq, Payload: p}, rpc.Reply{}
	if err = m.d.Call(req.ServiceMethod, &args, &reply); err != nil {
		return nil, err
	}
	resp := &Response{ServiceMethod: reply.ServiceMethod, Seq: reply.Seq, Data: req.Data}
	if err := oicq.Unmarshal(reply.Payload, resp.Data); err != nil {
		return nil, err
	}
	if err := resp.SetExtraData(resp.Data.TLVs); err != nil {
		return nil, err
	}
	p, _ = json.Marshal(resp)
	log.Trace("uin:%d seq:%d cmd:%s meta:%s", resp.Data.Uin, resp.Seq, resp.ServiceMethod, string(p))
	return resp, nil
}
