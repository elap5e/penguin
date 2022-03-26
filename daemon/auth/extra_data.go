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
	"encoding/hex"
	"log"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/auth/pb"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
)

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
			log.Printf("[D|ATHM|DUMP] t%x not parsed:\n%s", k, hex.Dump(v.MustGetValue().Bytes()))
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
			extraData.T547 = calcPow(extraData.T546)
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
