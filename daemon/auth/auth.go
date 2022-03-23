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
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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

	T119 []byte         `json:"t119,omitempty"`
	T150 []byte         `json:"t150,omitempty"`
	T161 []byte         `json:"t161,omitempty"`
	T174 []byte         `json:"t174,omitempty"`
	T17B []byte         `json:"t17b,omitempty"`
	T401 rpc.Key16Bytes `json:"t401,omitempty"`
	T402 []byte         `json:"t402,omitempty"`
	T403 []byte         `json:"t403,omitempty"`
	T542 []byte         `json:"t542,omitempty"`
	T546 []byte         `json:"t546,omitempty"`
	T547 []byte         `json:"t547,omitempty"`
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
	if v, ok := tlvs[0x000a].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		extraData.ErrorCode, _ = buf.ReadUint32()
		extraData.ErrorTitle, _ = buf.ReadStringL16V()
	}
	if v, ok := tlvs[0x0104].(*tlv.TLV); ok {
		extraData.SessionAuth = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0105].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		extraData.PictureSign, _ = buf.ReadBytesL16V()
		extraData.PictureData, _ = buf.ReadBytesL16V()
	}
	if v, ok := tlvs[0x0119].(*tlv.TLV); ok {
		extraData.T119 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0174].(*tlv.TLV); ok {
		extraData.T174 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x017b].(*tlv.TLV); ok {
		extraData.T17B = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0192].(*tlv.TLV); ok {
		extraData.CaptchaSign = string(v.MustGetValue().Bytes())
	}
	if v, ok := tlvs[0x0146].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		extraData.ErrorCode, _ = buf.ReadUint32()
		extraData.ErrorTitle, _ = buf.ReadStringL16V()
		extraData.ErrorMessage, _ = buf.ReadStringL16V()
	}
	if v, ok := tlvs[0x0150].(*tlv.TLV); ok {
		extraData.T150 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0161].(*tlv.TLV); ok {
		extraData.T161 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x017e].(*tlv.TLV); ok {
		extraData.Message = string(v.MustGetValue().Bytes())
	}
	if v, ok := tlvs[0x0178].(*tlv.TLV); ok {
		buf, _ := v.GetValue()
		_, _ = buf.ReadStringL16V()
		mobile, _ := buf.ReadStringL16V()
		extraData.SMSMobile = mobile
	}
	if v, ok := tlvs[0x0402].(*tlv.TLV); ok {
		extraData.T402 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0403].(*tlv.TLV); ok {
		extraData.T403 = v.MustGetValue().Bytes()
	}
	if v, ok := tlvs[0x0537].(*tlv.TLV); ok {
		extraData.Login, _ = v.MustGetValue().ReadBytesL16V()
	}
	if v, ok := tlvs[0x0546].(*tlv.TLV); ok {
		extraData.T546 = v.MustGetValue().Bytes()
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
	if req.Seq < 1 {
		req.Seq = m.c.GetNextSeq()
	}

	data := req.Data
	sess := m.c.GetSession(data.Uin)
	data.RandomKey = sess.RandomKey
	data.KeyVersion = sess.KeyVersion
	data.PublicKey = sess.PrivateKey.Public().Bytes()
	data.SharedSecret = sess.SharedSecret

	p, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("d.athm.request:\n%s", string(p))
	p, err := oicq.Marshal(data)
	if err != nil {
		return nil, err
	}

	args, reply := rpc.Args{
		Uin:     data.Uin,
		Seq:     req.Seq,
		Payload: p,
	}, rpc.Reply{}
	if err = m.c.Call(req.ServiceMethod, &args, &reply); err != nil {
		return nil, err
	}

	resp := &Response{
		ServiceMethod: reply.ServiceMethod,
		Seq:           reply.Seq,
		Data:          data,
	}
	if err := oicq.Unmarshal(reply.Payload, data); err != nil {
		return nil, err
	}
	if err := resp.SetExtraData(data.TLVs); err != nil {
		return nil, err
	}
	p, _ = json.MarshalIndent(resp, "", "  ")
	log.Printf("d.athm.response:\n%s", string(p))

	extraData := m.GetExtraData(data.Uin)
	switch data.Code {
	case 0x00:
		// success
		extraData.Login = resp.ExtraData.Login

		// decode t119
		tickets := m.c.GetTickets(data.Uin)
		key := [16]byte{}
		switch data.Type {
		default:
			log.Printf("unknown type: 0x%02x", data.Type)
			copy(key[:], tickets.A1.Key[:])
		case 0x0007: // AuthCheckSMSAndGetSessionTickets
			copy(key[:], tickets.A1.Key[:])
		case 0x0009: // signInWithPassword
			copy(key[:], tickets.A1.Key[:])
		case 0x000a: // signInWithoutPassword.A2
			copy(key[:], tickets.A2.Key[:])
		case 0x000b: // signInWithoutPassword.D2
			key = md5.Sum(tickets.D2.Key[:])
		case 0x0014: // AuthUnlockDevice
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
		extraData.T547 = resp.ExtraData.T546 // TODO: check
		if resp.ExtraData.CaptchaSign != "" {
			log.Println("verify captcha, url:", strings.ReplaceAll(resp.ExtraData.CaptchaSign, "https://", "https://elap5e.github.io/kits/"))
			// TODO: verify captcha
			fmt.Printf(">>> ")
			fmt.Scanln(&code)
			log.Printf("[%s]\n", code)
			return m.VerifyCaptcha(data.Uin, []byte(code))
		} else {
			log.Println("verify picture")
			// TODO: verify picture
			return m.VerifyCaptcha(data.Uin, []byte("code"), resp.ExtraData.PictureSign)
		}
	case 0xa0:
		// verify sms code
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.T17B = resp.ExtraData.T17B
		log.Println("verify sms code")
		// TODO: verify sms code
		return m.VerifySMSCode(data.Uin, []byte("code"))
	case 0xcc:
		// unlock device
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.T402 = resp.ExtraData.T402
		extraData.T403 = resp.ExtraData.T403
		extraData.T401 = md5.Sum(
			append(append(
				m.c.GetFakeSource(data.Uin).Device.GUID[:],
				sess.RandomPass[:]...),
				extraData.T402...),
		)
		log.Println("unlock device")
		return m.unlockDevice(data.Uin)
	case 0xef:
		// resend sms code
		m.c.SetSessionAuth(data.Uin, resp.ExtraData.SessionAuth)

		extraData.T174 = resp.ExtraData.T174
		extraData.T402 = resp.ExtraData.T402
		extraData.T403 = resp.ExtraData.T403
		extraData.T401 = md5.Sum(
			append(append(
				m.c.GetFakeSource(data.Uin).Device.GUID[:],
				sess.RandomPass[:]...),
				extraData.T402...),
		)
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
	case 0xed:
		log.Printf("invalid device, error:%s message:%s", resp.ExtraData.ErrorTitle, resp.ExtraData.ErrorMessage)
		return nil, fmt.Errorf("too many failures")
	default:
		log.Printf("unknown code: 0x%02x", data.Code)
	}
	return resp, nil
}

func checkUsername(username string) bool {
	uin, err := strconv.Atoi(username)
	if err != nil || uin < 10000 || uin > 4000000000 {
		return false
	}
	return true
}
