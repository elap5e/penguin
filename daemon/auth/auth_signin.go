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
	"strconv"

	"github.com/elap5e/penguin/daemon/constant"
	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/encoding/oicq"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/encoding/tlv/pb"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

// ACTION_WTLOGIN_GET_ST_WITH_PASSWD
// ACTION_WTLOGIN_GET_ST_WITHOUT_PASSWD
func (m *Manager) SignIn(username, password string) (*Response, error) {
	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	tickets := m.c.GetTickets(uin)
	if len(tickets.D2.Key) != 0 {
		return m.signInWithoutPassword(uin, false)
	} else if len(tickets.A2.Key) != 0 {
		return m.signInWithoutPassword(uin, true)
	}
	return m.signInWithPassword(uin, md5.Sum([]byte(password)))
}

func (m *Manager) signInWithPassword(uin int64, hash [16]byte) (*Response, error) {
	username := strconv.FormatInt(uin, 10)
	fake, sess, seq, serverTime := m.c.GetFakeSource(uin), m.c.GetSession(uin), m.c.GetNextSeq(), m.c.GetServerTime()
	extraData := m.GetExtraData(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0018] = tlv.NewT18(constant.DstAppID, 0, uin, 0)
	tlvs[0x0001] = tlv.NewT1(uin, fake.Device.IPv4Address, serverTime)
	a1 := m.c.GetTickets(uin).A1
	if len(a1.Sig) == 0 {
		tlvs[0x0106] = tlv.NewT106(constant.DstAppID, constant.OpenAppID, 0, uin, serverTime, fake.Device.IPv4Address, true, hash, 0, username, a1.Key, true, fake.Device.GUID[:], 1, fake.App.SSOVer)
	} else {
		tlvs[0x0106] = tlv.NewTLV(0x0106, 0x0000, bytes.NewBuffer(a1.Sig))
	}
	tlvs[0x0116] = tlv.NewT116(fake.App.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0100] = tlv.NewT100(constant.DstAppID, constant.OpenAppID, 0, constant.MainSigMap, fake.App.SSOVer)
	tlvs[0x0107] = tlv.NewT107(0, 0, 0, 1)
	if len(sess.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(sess.KSID)
	}
	if len(sess.Auth) != 0 {
		tlvs[0x0104] = tlv.NewT104(sess.Auth)
	}
	tlvs[0x0142] = tlv.NewT142([]byte(fake.App.PkgName))
	if !checkUsername(username) {
		tlvs[0x0112] = tlv.NewT112([]byte(username))
	}
	tlvs[0x0144] = tlv.NewT144(a1.Key,
		tlv.NewT109(md5.Sum([]byte(fake.Device.OS.BuildID))),
		tlv.NewT52D(&pb.DeviceReport{
			Bootloader:   []byte(fake.Device.Bootloader),
			ProcVersion:  []byte(fake.Device.ProcVersion),
			Codename:     []byte(fake.Device.Codename),
			Incremental:  []byte(fake.Device.Incremental),
			Fingerprint:  []byte(fake.Device.Fingerprint),
			BootId:       []byte(fake.Device.BootID),
			AndroidId:    []byte(fake.Device.OS.BuildID),
			Baseband:     []byte(fake.Device.Baseband),
			InnerVersion: []byte(fake.Device.InnerVersion),
		}),
		tlv.NewT124(
			[]byte(fake.Device.OS.Type),
			[]byte(fake.Device.OS.Version),
			fake.Device.OS.NetworkType,
			fake.Device.SIMOPName,
			nil,
			fake.Device.APNName,
		),
		tlv.NewT128(
			fake.Device.IsGUIDFileNil,
			fake.Device.IsGUIDGenSucc,
			fake.Device.IsGUIDChanged,
			fake.Device.GUIDFlag,
			[]byte(fake.Device.OS.BuildModel),
			fake.Device.GUID[:],
			fake.Device.OS.BuildBrand,
		),
		tlv.NewT16E([]byte(fake.Device.OS.BuildModel)),
	)
	tlvs[0x0145] = tlv.NewT145(fake.Device.GUID)
	tlvs[0x0147] = tlv.NewT147(constant.DstAppID, []byte(fake.App.VerName), fake.App.SigHash)
	if fake.App.MiscBitMap&0x80 != 0 {
		tlvs[0x0166] = tlv.NewT166(fake.App.ImageType)
	}
	if len(extraData.T16A) != 0 {
		tlvs[0x016a] = tlv.NewT16A(extraData.T16A)
	}
	tlvs[0x0154] = tlv.NewT154(seq)
	tlvs[0x0141] = tlv.NewT141(fake.Device.SIMOPName, fake.Device.OS.NetworkType, fake.Device.APNName)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0511] = tlv.NewT511(constant.AppDomains)
	if len(extraData.T172) != 0 {
		tlvs[0x0172] = tlv.NewT172(extraData.T172)
	}
	// TODO: LoginType == 3
	// tlvs[0x0185] = tlv.NewT185(0x01)
	// TODO: code2d
	// tlvs[0x0400] = tlv.NewT400(
	// 	h.hashedGUID,
	// 	req.GetUin(),
	// 	util.BytesToSTBytes(fake.Device.GUID),
	// 	h.randomPassword,
	// 	req.DstAppID,
	// 	req.SubDstAppID,
	// 	h.randomSeed,
	// )
	tlvs[0x0187] = tlv.NewT187(md5.Sum([]byte(fake.Device.MACAddress)))
	tlvs[0x0188] = tlv.NewT188(md5.Sum([]byte(fake.Device.OS.BuildID)))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(fake.Device.IMSI)))
	if fake.App.CanCaptcha {
		tlvs[0x0191] = tlv.NewT191(0x82)
	} else {
		tlvs[0x0191] = tlv.NewT191(0x00)
	}
	// DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum([]byte(fake.Device.BSSIDAddress)), []byte(fake.Device.SSIDAddress))
	tlvs[0x0177] = tlv.NewT177(fake.App.BuildAt, fake.App.SDKVer)
	tlvs[0x0516] = tlv.NewTLV(0x0516, 0x0004, bytes.NewBuffer([]byte{0, 0, 0, 0}))       // SourceType
	tlvs[0x0521] = tlv.NewTLV(0x0521, 0x0006, bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0})) // ProductType
	if len(extraData.Login) != 0 {
		buf := bytes.NewBuffer([]byte{})
		buf.WriteUint16(0x0001)
		tlv.NewTLV(0x0536, 0x0002, bytes.NewBuffer(extraData.Login)).WriteTo(buf)
		tlvs[0x0525] = tlv.NewTLV(0x0525, 0x0000, buf)
	}
	// TODO: ???
	// tlvs[0x0529] = tlv.NewT529()
	// TODO: code2d
	// if len(h.tgtQR) != 0 {
	// 	tlvs[0x0318] = tlv.NewTLV(0x0318, 0x0000, bytes.NewBuffer(h.tgtQR))
	// }
	// DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(uin, fake.Device.GUID, fake.SDKVer, 9)
	// DISABLED: tgtgt qimei
	// tlvs[0x0545] = tlv.NewT545(md5.Sum([]byte("qimei")))
	// DISABLED: nativeGetTestData
	// tlvs[0x0548] = tlv.NewT548([]byte("nativeGetTestData"))
	return m.requestSignIn(seq, uin, 9, tlvs)
}

func (m *Manager) requestSignIn(seq int32, uin int64, typ uint16, tlvs map[uint16]tlv.Codec) (*Response, error) {
	return m.request(&Request{
		ServiceMethod: service.MethodAuthSignIn,
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

func (m *Manager) signInWithoutPassword(uin int64, changeD2 bool) (*Response, error) {
	username := strconv.FormatInt(uin, 10)
	tickets := m.c.GetTickets(uin)
	fake, sess, seq := m.c.GetFakeSource(uin), m.c.GetSession(uin), m.c.GetNextSeq()
	extraData := m.GetExtraData(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0100] = tlv.NewT100(constant.DstAppID, constant.OpenAppID, 0, constant.MainSigMap, fake.App.SSOVer)
	tlvs[0x010a] = tlv.NewT10A(tickets.A2.Sig)
	tlvs[0x0116] = tlv.NewT116(fake.App.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	if len(sess.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(sess.KSID)
	}
	key := [16]byte{}
	if !changeD2 {
		copy(key[:], tickets.A2.Key[:])
	} else {
		key = md5.Sum(tickets.D2.Key[:])
	}
	tlvs[0x0144] = tlv.NewT144(key,
		tlv.NewT109(md5.Sum([]byte(fake.Device.OS.BuildID))),
		tlv.NewT52D(&pb.DeviceReport{
			Bootloader:   []byte(fake.Device.Bootloader),
			ProcVersion:  []byte(fake.Device.ProcVersion),
			Codename:     []byte(fake.Device.Codename),
			Incremental:  []byte(fake.Device.Incremental),
			Fingerprint:  []byte(fake.Device.Fingerprint),
			BootId:       []byte(fake.Device.BootID),
			AndroidId:    []byte(fake.Device.OS.BuildID),
			Baseband:     []byte(fake.Device.Baseband),
			InnerVersion: []byte(fake.Device.InnerVersion),
		}),
		tlv.NewT124(
			[]byte(fake.Device.OS.Type),
			[]byte(fake.Device.OS.Version),
			fake.Device.OS.NetworkType,
			fake.Device.SIMOPName,
			nil,
			fake.Device.APNName,
		),
		tlv.NewT128(
			fake.Device.IsGUIDFileNil,
			fake.Device.IsGUIDGenSucc,
			fake.Device.IsGUIDChanged,
			fake.Device.GUIDFlag,
			[]byte(fake.Device.OS.BuildModel),
			fake.Device.GUID[:],
			fake.Device.OS.BuildBrand,
		),
		tlv.NewT16E([]byte(fake.Device.OS.BuildModel)),
	)
	if !checkUsername(username) {
		tlvs[0x0112] = tlv.NewT112([]byte(username))
	}
	if !changeD2 {
		tlvs[0x0145] = tlv.NewT145(fake.Device.GUID)
	} else {
		tlvs[0x0143] = tlv.NewT143(tickets.D2.Sig)
	}
	tlvs[0x0142] = tlv.NewT142([]byte(fake.App.PkgName))
	tlvs[0x0154] = tlv.NewT154(seq)
	tlvs[0x0018] = tlv.NewT18(constant.DstAppID, 0, uin, 0)
	tlvs[0x0141] = tlv.NewT141(fake.Device.SIMOPName, fake.Device.OS.NetworkType, fake.Device.APNName)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0511] = tlv.NewT511(constant.AppDomains)
	tlvs[0x0147] = tlv.NewT147(constant.DstAppID, []byte(fake.App.VerName), fake.App.SigHash)
	if len(extraData.T172) != 0 {
		tlvs[0x0172] = tlv.NewT172(extraData.T172)
	}
	tlvs[0x0177] = tlv.NewT177(fake.App.BuildAt, fake.App.SDKVer)
	tlvs[0x0187] = tlv.NewT187(md5.Sum([]byte(fake.Device.MACAddress)))
	tlvs[0x0188] = tlv.NewT188(md5.Sum([]byte(fake.Device.OS.BuildID)))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(fake.Device.IMSI)))
	// DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum([]byte(fake.Device.BSSIDAddress)), []byte(fake.Device.SSIDAddress))
	// DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(uin, fake.Device.GUID, fake.SDKVer, 10)
	if !changeD2 {
		return m.requestSignIn(seq, uin, 10, tlvs)
	}
	return m.requestSignIn(seq, uin, 11, tlvs)
}

func (m *Manager) requestSignInA2(seq int32, uin int64, typ uint16, tlvs map[uint16]tlv.Codec) (*Response, error) {
	return m.request(&Request{
		ServiceMethod: service.MethodAuthSignInA2,
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
