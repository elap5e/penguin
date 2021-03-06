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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/elap5e/penguin/daemon/constant"
	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/encoding/oicq"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/encoding/tlv/pb"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

const (
	LoginTypeUin = uint32(1)
	LoginTypeSMS = uint32(3)
)

func (m *Manager) SignIn(username, password string) (*Response, error) {
	if strings.HasPrefix(username, "+") {
		if strings.HasPrefix(username, "+86") {
			username = username[3:]
		} else {
			username = "00" + username[1:]
		}
		token, err := m.verifySignInWithCodeCaptach(username)
		if err != nil {
			return nil, err
		}
		return m.signInWithCode(username, token)
	} else if !checkUsername(username) {
		re, err := regexp.Compile(`^([0-9]{5,10})@qq\.com$`)
		if err != nil {
			return nil, err
		}
		matched := re.FindStringSubmatch(username)
		if matched == nil || len(matched) != 2 {
			return nil, fmt.Errorf("not a valid username")
		}
		username = matched[1]
	}
	_ = m.SetAccount(username, password)
	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	tickets := m.c.GetTickets(uin)
	if tickets.A2.Valid() {
		return m.SignInUpdateToken(username, true)
	} else if tickets.D2.Valid() {
		return m.SignInUpdateToken(username)
	}
	return m.SignInCreateToken(username, uin)
}

func (m *Manager) SignInCreateToken(username string, uin int64) (*Response, error) {
	account := m.GetAccount(username)
	if account.Generate {
		return nil, fmt.Errorf("account %s is generated", username)
	}
	return m.signInWithPassword(account.Username, account.Password, uin, 0, LoginTypeUin)
}

func (m *Manager) SignInUpdateToken(username string, changeD2 ...bool) (*Response, error) {
	if len(changeD2) > 0 && changeD2[0] {
		return m.signInWithoutPassword(username, true)
	}
	return m.signInWithoutPassword(username, false)
}

// ACTION_WTLOGIN_GET_ST_WITH_PASSWD
// ACTION_WTLOGIN_GET_ST_VIA_SMS_VERIFY_LOGIN
func (m *Manager) signInWithPassword(username string, hash [16]byte, uin, salt int64, loginType uint32) (*Response, error) {
	fake, session, tickets, seq, serverTime := m.c.GetFakeSource(uin), m.c.GetSession(uin), m.c.GetTickets(uin), m.c.GetNextSeq(), m.c.GetServerTime()
	extraData := m.GetExtraData(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0018] = tlv.NewT18(constant.DstAppID, 0, uin, 0)
	tlvs[0x0001] = tlv.NewT1(uin, fake.Device.IPv4Address, serverTime)
	a1 := m.c.GetTickets(uin).A1
	if len(a1.Sig) == 0 {
		tlvs[0x0106] = tlv.NewT106(constant.DstAppID, constant.OpenAppID, 0, uin, serverTime, fake.Device.IPv4Address, true, hash, salt, username, a1.Key, true, fake.Device.GUID[:], loginType, fake.SDK.SSOVersion)
	} else {
		tlvs[0x0106] = tlv.NewTLV(0x0106, 0x0000, bytes.NewBuffer(a1.Sig))
	}
	tlvs[0x0116] = tlv.NewT116(constant.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0100] = tlv.NewT100(constant.DstAppID, constant.OpenAppID, 0, constant.MainSigMap, fake.SDK.SSOVersion)
	tlvs[0x0107] = tlv.NewT107(0, 0, 0, 1)
	if len(tickets.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(tickets.KSID)
	}
	if len(session.Auth) != 0 {
		tlvs[0x0104] = tlv.NewT104(session.Auth)
	}
	tlvs[0x0142] = tlv.NewT142([]byte(fake.App.Package))
	if !checkUsername(username) {
		tlvs[0x0112] = tlv.NewT112([]byte(username))
	}
	tlvs[0x0144] = tlv.NewT144(a1.Key,
		tlv.NewT109(md5.Sum([]byte(fake.Device.OS.BuildID))),
		tlv.NewT52D(&pb.DeviceReport{
			Bootloader:  []byte(fake.Device.Bootloader),
			Version:     []byte(fake.Device.ProcVersion),
			Codename:    []byte(fake.Device.Codename),
			Incremental: []byte(fake.Device.Incremental),
			Fingerprint: []byte(fake.Device.Fingerprint),
			BootId:      []byte(fake.Device.BootID),
			AndroidId:   []byte(fake.Device.OS.BuildID),
			Baseband:    []byte(fake.Device.Baseband),
			InnerVer:    []byte(fake.Device.InnerVersion),
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
	tlvs[0x0147] = tlv.NewT147(constant.DstAppID, []byte(fake.App.Version), fake.App.SigHash)
	if constant.MiscBitMap&0x80 != 0 {
		tlvs[0x0166] = tlv.NewT166(fake.SDK.ImageType)
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
	if loginType == LoginTypeSMS {
		tlvs[0x0185] = tlv.NewT185(1)
	}
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
	if fake.SDK.CanCaptcha {
		tlvs[0x0191] = tlv.NewT191(0x82)
	} else {
		tlvs[0x0191] = tlv.NewT191(0x00)
	}
	// DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum([]byte(fake.Device.BSSIDAddress)), []byte(fake.Device.SSIDAddress))
	tlvs[0x0177] = tlv.NewT177(fake.SDK.BuildAt, fake.SDK.Version)
	tlvs[0x0516] = tlv.NewTLV(0x0516, 0x0004, bytes.NewBuffer([]byte{0, 0, 0, 0}))       // SourceType
	tlvs[0x0521] = tlv.NewTLV(0x0521, 0x0006, bytes.NewBuffer([]byte{0, 0, 0, 0, 0, 0})) // ProductType
	if len(extraData.Login) != 0 {
		buf := bytes.NewBuffer([]byte{})
		buf.WriteUint16(1)
		tlv.NewTLV(0x0536, 0x0002, bytes.NewBuffer(extraData.Login)).WriteTo(buf)
		tlvs[0x0525] = tlv.NewTLV(0x0525, 0, buf)
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
			EncryptMethod: oicq.EncryptMethodECDH0x87,
			Type:          typ,
			TLVs:          tlvs,
		},
	})
}

// ACTION_WTLOGIN_GET_ST_WITHOUT_PASSWD
// ACTION_WTLOGIN_GET_OPEN_KEY_WITHOUT_PASSWD
func (m *Manager) signInWithoutPassword(username string, changeD2 bool) (*Response, error) {
	uin, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return nil, err
	}
	tickets := m.c.GetTickets(uin)
	fake, tickets, seq := m.c.GetFakeSource(uin), m.c.GetTickets(uin), m.c.GetNextSeq()
	extraData := m.GetExtraData(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0100] = tlv.NewT100(constant.DstAppID, constant.OpenAppID, 0, constant.MainSigMap, fake.SDK.SSOVersion)
	tlvs[0x010a] = tlv.NewT10A(tickets.A2.Sig)
	tlvs[0x0116] = tlv.NewT116(constant.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	if len(tickets.KSID) != 0 {
		tlvs[0x0108] = tlv.NewT108(tickets.KSID)
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
			Bootloader:  []byte(fake.Device.Bootloader),
			Version:     []byte(fake.Device.ProcVersion),
			Codename:    []byte(fake.Device.Codename),
			Incremental: []byte(fake.Device.Incremental),
			Fingerprint: []byte(fake.Device.Fingerprint),
			BootId:      []byte(fake.Device.BootID),
			AndroidId:   []byte(fake.Device.OS.BuildID),
			Baseband:    []byte(fake.Device.Baseband),
			InnerVer:    []byte(fake.Device.InnerVersion),
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
	tlvs[0x0142] = tlv.NewT142([]byte(fake.App.Package))
	tlvs[0x0154] = tlv.NewT154(seq)
	tlvs[0x0018] = tlv.NewT18(constant.DstAppID, 0, uin, 0)
	tlvs[0x0141] = tlv.NewT141(fake.Device.SIMOPName, fake.Device.OS.NetworkType, fake.Device.APNName)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0511] = tlv.NewT511(constant.AppDomains)
	tlvs[0x0147] = tlv.NewT147(constant.DstAppID, []byte(fake.App.Version), fake.App.SigHash)
	if len(extraData.T172) != 0 {
		tlvs[0x0172] = tlv.NewT172(extraData.T172)
	}
	tlvs[0x0177] = tlv.NewT177(fake.SDK.BuildAt, fake.SDK.Version)
	tlvs[0x0187] = tlv.NewT187(md5.Sum([]byte(fake.Device.MACAddress)))
	tlvs[0x0188] = tlv.NewT188(md5.Sum([]byte(fake.Device.OS.BuildID)))
	tlvs[0x0194] = tlv.NewT194(md5.Sum([]byte(fake.Device.IMSI)))
	// DISABLED: SetNeedForPayToken
	// tlvs[0x0201] = tlv.NewT201(nil, nil, []byte("qq"), nil)
	tlvs[0x0202] = tlv.NewT202(md5.Sum([]byte(fake.Device.BSSIDAddress)), []byte(fake.Device.SSIDAddress))
	// DISABLED: tgt
	// tlvs[0x0544] = tlv.NewT544(uin, fake.Device.GUID, fake.SDKVer, 10)
	if !changeD2 {
		return m.requestSignInA2(seq, uin, 10, tlvs)
	}
	return m.requestSignInA2(seq, uin, 11, tlvs)
}

func (m *Manager) requestSignInA2(seq int32, uin int64, typ uint16, tlvs map[uint16]tlv.Codec) (*Response, error) {
	return m.request(&Request{
		ServiceMethod: service.MethodAuthSignInA2,
		Seq:           seq,
		Data: &oicq.Data{
			Version:       0x1f41,
			ServiceMethod: 0x0810,
			Uin:           uin,
			EncryptMethod: oicq.EncryptMethodECDH0x87,
			Type:          typ,
			TLVs:          tlvs,
		},
	})
}
