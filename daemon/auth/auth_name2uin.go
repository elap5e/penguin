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
	"github.com/elap5e/penguin/pkg/encoding/oicq"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/encoding/tlv/pb"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

// ACTION_WTLOGIN_REFRESH_SMS_DATA
func (m *Manager) name2Uin(username string) (*Response, error) {
	fake, sess, seq := m.c.GetFakeSource(0), m.c.GetSession(0), m.c.GetNextSeq()
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0100] = tlv.NewT100(constant.DstAppID, constant.OpenAppID, 0, constant.MainSigMap, fake.App.SSOVer)
	tlvs[0x0112] = tlv.NewT112([]byte(username))
	tlvs[0x0107] = tlv.NewT107(0, 0, 0, 1)
	tlvs[0x0154] = tlv.NewT154(seq)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	if len(sess.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(sess.KSID)
	}
	tlvs[0x0521] = tlv.NewTLV(0x0521, 6, bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0})) // ProductType
	tlvs[0x0124] = tlv.NewT124(
		[]byte(fake.Device.OS.Type),
		[]byte(fake.Device.OS.Version),
		fake.Device.OS.NetworkType,
		fake.Device.SIMOPName,
		nil,
		fake.Device.APNName,
	)
	tlvs[0x0128] = tlv.NewT128(
		fake.Device.IsGUIDFileNil,
		fake.Device.IsGUIDGenSucc,
		fake.Device.IsGUIDChanged,
		fake.Device.GUIDFlag,
		[]byte(fake.Device.OS.BuildModel),
		fake.Device.GUID[:],
		fake.Device.OS.BuildBrand,
	)
	tlvs[0x0116] = tlv.NewT116(fake.App.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0191] = tlv.NewT191(0x82)
	tlvs[0x011b] = tlv.NewTLV(0x011b, 1, bytes.NewBuffer([]byte{2}))
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
	// DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))
	return m.requestName2Uin(seq, 0, 4, tlvs)
}

func (m *Manager) requestName2Uin(seq int32, uin int64, typ uint16, tlvs map[uint16]tlv.Codec) (*Response, error) {
	return m.request(&Request{
		ServiceMethod: service.MethodAuthName2Uin,
		Seq:           seq,
		Data: &oicq.Data{
			Version:       0x1f41,
			ServiceMethod: 0x0810,
			Uin:           uin,
			EncryptMethod: oicq.EncryptMethodECDH,
			Type:          typ,
			TLVs:          tlvs,
		},
	})
}
