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
	"github.com/elap5e/penguin/daemon/constant"
	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
)

func (m *Manager) resendSMSCode(uin int64) (*Response, error) {
	extraData := m.GetExtraData(uin)
	sess := m.c.GetSession(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0104] = tlv.NewT104(sess.Auth)
	tlvs[0x0116] = tlv.NewT116(m.c.GetFake(uin).MiscBitmap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0174] = tlv.NewT174(extraData.T174)
	tlvs[0x017a] = tlv.NewT17A(constant.SMSAppID)
	tlvs[0x0197] = tlv.NewTLV(0x0197, 0x0000, bytes.NewBuffer([]byte{0}))
	tlvs[0x0542] = tlv.NewT542(extraData.T542)
	return m.requestSignIn(0, uin, 0x0008, tlvs)
}

func (m *Manager) VerifyCaptcha(uin int64, code []byte, sign ...[]byte) (*Response, error) {
	extraData := m.GetExtraData(uin)
	sess := m.c.GetSession(uin)
	tlvs := make(map[uint16]tlv.Codec)
	if len(sign) > 0 {
		tlvs[0x0193] = tlv.NewT193(code)
	} else {
		tlvs[0x0002] = tlv.NewT2(code, sign[0])
	}
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0104] = tlv.NewT104(sess.Auth)
	tlvs[0x0116] = tlv.NewT116(m.c.GetFake(uin).MiscBitmap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0547] = tlv.NewT547(extraData.T547)
	return m.requestSignIn(0, uin, 0x0002, tlvs)
}

func (m *Manager) VerifySMSCode(uin int64, code []byte) (*Response, error) {
	extraData := m.GetExtraData(uin)
	sess := m.c.GetSession(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0104] = tlv.NewT104(sess.Auth)
	tlvs[0x0116] = tlv.NewT116(m.c.GetFake(uin).MiscBitmap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0174] = tlv.NewT174(extraData.T174)
	tlvs[0x017c] = tlv.NewT17C(code)
	tlvs[0x0401] = tlv.NewT401(extraData.T401)
	tlvs[0x0197] = tlv.NewTLV(0x0197, 0x0000, bytes.NewBuffer([]byte{0}))
	tlvs[0x0542] = tlv.NewT542(extraData.T542)
	return m.requestSignIn(0, uin, 0x0007, tlvs)
}
