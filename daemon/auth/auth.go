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
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/auth/pb"
	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
	"github.com/elap5e/penguin/pkg/encoding/oicq"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Manager struct {
	ctx context.Context

	c rpc.Client

	mapExtarData map[int64]*ExtraData
}

func NewManager(ctx context.Context, c rpc.Client) *Manager {
	return &Manager{
		ctx:          ctx,
		c:            c,
		mapExtarData: make(map[int64]*ExtraData),
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

func (resp *Response) SetExtraData(tlvs map[uint16]tlv.Codec) error {
	extraData := resp.ExtraData
	if extraData == nil {
		resp.ExtraData = &ExtraData{}
		extraData = resp.ExtraData
	}
	for k, v := range tlvs {
		v := v.(*tlv.TLV)
		switch k {
		default:
			log.Printf("t%x not parsed:\n%s\n", k, hex.Dump(v.MustGetValue().Bytes()))
		case 0x000a:
			buf, _ := v.GetValue()
			extraData.ErrorCode, _ = buf.ReadUint32()
			extraData.ErrorTitle, _ = buf.ReadStringL16V()
		case 0x0104:
			extraData.SessionAuth = v.MustGetValue().Bytes()
		case 0x0105:
			buf, _ := v.GetValue()
			extraData.PictureSign, _ = buf.ReadBytesL16V()
			extraData.PictureData, _ = buf.ReadBytesL16V()
		case 0x0119:
			extraData.T119 = v.MustGetValue().Bytes()
		case 0x0126:
			buf, _ := v.GetValue()
			_, _ = buf.ReadUint16()
			extraData.SignInCodeSign, _ = buf.ReadBytesL16V()
		case 0x0174:
			extraData.T174 = v.MustGetValue().Bytes()
		case 0x017b:
			extraData.T17B = v.MustGetValue().Bytes()
		case 0x0182:
			extraData.T182 = v.MustGetValue().Bytes()
		case 0x0183:
			extraData.Salt, _ = v.MustGetValue().ReadUint64()
		case 0x0192:
			extraData.CaptchaSign = string(v.MustGetValue().Bytes())
		case 0x0146:
			buf, _ := v.GetValue()
			extraData.ErrorCode, _ = buf.ReadUint32()
			extraData.ErrorTitle, _ = buf.ReadStringL16V()
			extraData.ErrorMessage, _ = buf.ReadStringL16V()
		case 0x0150:
			extraData.T150 = v.MustGetValue().Bytes()
		case 0x0161:
			extraData.T161 = v.MustGetValue().Bytes()
		case 0x017e:
			extraData.Message = string(v.MustGetValue().Bytes())
		case 0x0178:
			buf, _ := v.GetValue()
			_, _ = buf.ReadStringL16V()
			mobile, _ := buf.ReadStringL16V()
			extraData.SMSMobile = mobile
		case 0x0402:
			extraData.T402 = v.MustGetValue().Bytes()
		case 0x0403:
			extraData.T403 = v.MustGetValue().Bytes()
		case 0x0537:
			extraData.Login, _ = v.MustGetValue().ReadBytesL16V()
		case 0x0543:
			extraData.ThirdPartLogin = &pb.ThirdPartLogin_RspBody{}
			_ = proto.Unmarshal(v.MustGetValue().Bytes(), extraData.ThirdPartLogin)
		case 0x0546:
			extraData.T546 = v.MustGetValue().Bytes()
		}
	}
	return nil
}

func (m *Manager) GetExtraData(uin int64) *ExtraData {
	extraData := m.mapExtarData[uin]
	if extraData == nil {
		m.mapExtarData[uin] = &ExtraData{}
		extraData = m.mapExtarData[uin]
	}
	return extraData
}

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
			log.Println(login.UnbindWording)

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
			log.Printf("unknown type: 0x%02x", data.Type)
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
		if v, ok := tlvs[0x0108]; ok {
			m.c.SetSessionKSID(data.Uin, v.(*tlv.TLV).MustGetValue().Bytes())
		}
		return resp, nil
	case 0x02:
		// verify captcha or picture
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		var code string
		extraData.T546 = resp.ExtraData.T546 // TODO: check
		if resp.ExtraData.CaptchaSign != "" {
			l, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return nil, err
			}
			addr := l.Addr().(*net.TCPAddr).String()
			log.Println("verify captcha, url:", strings.ReplaceAll(
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
			log.Println("verify picture")
			fmt.Printf(">>> ")
			fmt.Scanln(&code)
			return m.VerifyCaptcha(data.Uin, []byte(code), resp.ExtraData.PictureSign)
		}
	case 0xa0:
		// verify sms code
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.T17B = resp.ExtraData.T17B
		log.Println("verify sms code")
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
		log.Println("verify signin code")
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
		log.Println("unlock device")
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
		log.Println("resend sms code, mobile:", resp.ExtraData.SMSMobile)
		return m.resendSMSCode(data.Uin)
	case 0x01:
		log.Printf("invalid login, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("invalid password")
	case 0x06:
		log.Printf("not implement, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x09:
		log.Printf("invalid service, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x0a:
		log.Printf("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("service temporarily unavailable")
	case 0x10:
		log.Printf("session expired, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x28:
		log.Printf("protection mode, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
	case 0x9a:
		log.Printf("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("service temporarily unavailable")
	case 0xa1:
		log.Printf("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("too many sms verify requests")
	case 0xa2:
		log.Printf("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("frequent sms verify requests")
	case 0xa4:
		log.Printf("error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("bad requests")
	case 0xd5:
		log.Printf("phone number not valid, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("phone number not valid")
	case 0xdb:
		log.Printf("phone number not registered, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("phone number not registered")
	case 0xeb:
		log.Printf("version too low, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("too many failures")
	case 0xed:
		log.Printf("invalid device, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("too many failures")
	default:
		log.Printf("unknown code: 0x%02x", data.Code)
	}
	return resp, nil
}

func (m *Manager) call(req *Request) (*Response, error) {
	p, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("d.athm.request:\n%s", string(p))
	p, err := oicq.Marshal(req.Data)
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{Uin: req.Data.Uin, Seq: req.Seq, Payload: p}, rpc.Reply{}
	if err = m.c.Call(req.ServiceMethod, &args, &reply); err != nil {
		return nil, err
	}
	resp := &Response{ServiceMethod: reply.ServiceMethod, Seq: reply.Seq, Data: req.Data}
	if err := oicq.Unmarshal(reply.Payload, resp.Data); err != nil {
		return nil, err
	}
	if err := resp.SetExtraData(resp.Data.TLVs); err != nil {
		return nil, err
	}
	p, _ = json.MarshalIndent(resp, "", "  ")
	log.Printf("d.athm.response:\n%s", string(p))
	return resp, nil
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
