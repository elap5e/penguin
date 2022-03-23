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

package tlv

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"net"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
)

type T106 struct {
	tlv      *TLV
	appID    uint64
	subAppID uint64
	acVer    uint32
	uin      int64
	svrTime  int64
	ip       net.IP
	i2       bool
	hash     [16]byte
	salt     int64
	username string
	a1Kay    [16]byte
	haveGUID bool
	guid     []byte
	typ      uint32

	ssoVer uint32
}

func NewT106(appID, subAppID uint64, acVer uint32, uin int64, svrTime int64, ip net.IP, i2 bool, hash [16]byte, salt int64, username string, a1Kay [16]byte, haveGUID bool, guid []byte, typ, ssoVer uint32) *T106 {
	return &T106{
		tlv:      NewTLV(0x0106, 0x0000, nil),
		appID:    appID,
		subAppID: subAppID,
		acVer:    acVer,
		uin:      uin,
		svrTime:  svrTime,
		ip:       ip,
		i2:       i2,
		hash:     hash,
		salt:     salt,
		username: username,
		a1Kay:    a1Kay,
		haveGUID: haveGUID,
		guid:     guid,
		typ:      typ,

		ssoVer: ssoVer,
	}
}

func (t *T106) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}

func (t *T106) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(4)
	v.WriteUint32(rand.Uint32())
	v.WriteUint32(t.ssoVer)
	v.WriteUint32(uint32(t.appID))
	v.WriteUint32(t.acVer)
	if t.uin != 0 {
		v.WriteUint64(uint64(t.uin))
	} else {
		v.WriteUint64(uint64(t.salt))
	}
	v.WriteUint32(uint32(t.svrTime))
	v.Write(t.ip.To4())
	v.WriteBool(t.i2)
	v.Write(t.hash[:])
	v.Write(t.a1Kay[:])
	v.WriteUint32(0)
	v.WriteBool(t.haveGUID)
	if len(t.guid) == 0 {
		v.WriteUint64(rand.Uint64())
		v.WriteUint64(rand.Uint64())
	} else {
		v.Write(t.guid)
	}
	v.WriteUint32(uint32(t.subAppID))
	v.WriteUint32(t.typ)
	v.WriteStringL16V(t.username)

	key := append(t.hash[:], make([]byte, 8)...)
	if t.uin != 0 {
		binary.BigEndian.PutUint64(key[16:], uint64(t.uin))
	} else {
		binary.BigEndian.PutUint64(key[16:], uint64(t.salt))
	}
	t.tlv.SetValue(bytes.NewBuffer(tea.NewCipher(md5.Sum(key)).Encrypt(v.Bytes())))
	return t.tlv.WriteTo(b)
}
