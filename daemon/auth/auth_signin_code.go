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
	"log"

	"github.com/elap5e/penguin/daemon/constant"
	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/encoding/tlv/pb"
)

func (m *Manager) signInWithCode(username string, token []byte) (*Response, error) {
	fake, tickets, seq := m.c.GetFakeSource(0), m.c.GetTickets(0), m.c.GetNextSeq()
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0100] = tlv.NewT100(constant.DstAppID, constant.OpenAppID, 0, constant.MainSigMap, fake.App.SSOVer)
	if len(tickets.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(tickets.KSID)
	}
	tlvs[0x0109] = tlv.NewT109(md5.Sum([]byte(fake.Device.OS.BuildID)))
	tlvs[0x052d] = tlv.NewT52D(&pb.DeviceReport{
		Bootloader:  []byte(fake.Device.Bootloader),
		Version:     []byte(fake.Device.ProcVersion),
		Codename:    []byte(fake.Device.Codename),
		Incremental: []byte(fake.Device.Incremental),
		Fingerprint: []byte(fake.Device.Fingerprint),
		BootId:      []byte(fake.Device.BootID),
		AndroidId:   []byte(fake.Device.OS.BuildID),
		Baseband:    []byte(fake.Device.Baseband),
		InnerVer:    []byte(fake.Device.InnerVersion),
	})
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0142] = tlv.NewT142([]byte(fake.App.PkgName))
	tlvs[0x0145] = tlv.NewT145(fake.Device.GUID)
	tlvs[0x0154] = tlv.NewT154(seq)
	tlvs[0x0112] = tlv.NewT112([]byte(username))
	tlvs[0x0116] = tlv.NewT116(constant.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	// DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))
	tlvs[0x0542] = tlv.NewT542(token)
	return m.requestSignIn(seq, 0, 17, tlvs)
}

func (m *Manager) VerifySignInCode(code []byte) (*Response, error) {
	extraData, session := m.GetExtraData(0), m.c.GetSession(0)
	password := randomPassword()
	log.Println(password)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0104] = tlv.NewT104(session.Auth)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0127] = tlv.NewT127(code, extraData.SignInCodeSign)
	tlvs[0x0184] = tlv.NewT184(extraData.Salt, password)
	tlvs[0x0116] = tlv.NewT116(constant.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	// DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))
	return m.requestSignIn(0, 0, 18, tlvs)
}

// ACTION_WTLOGIN_REFRESH_SMS_VERIFY_LOGIN_CODE
func (m *Manager) resendSignInCode(uin int64) (*Response, error) {
	session := m.c.GetSession(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0104] = tlv.NewT104(session.Auth)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0116] = tlv.NewT116(constant.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0521] = tlv.NewTLV(0x0521, 6, bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0})) // ProductType
	return m.requestSignIn(0, uin, 19, tlvs)
}
