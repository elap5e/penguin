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
	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
)

type T400 struct {
	tlv      *TLV
	key      [16]byte
	uin      uint64
	guid     [16]byte
	dpwd     [16]byte
	appID    uint64
	subAppID uint64
	randSeed []byte

	serverTime int64
}

func NewT400(key [16]byte, uin uint64, guid, dpwd [16]byte, appID, subAppID uint64, randSeed []byte, serverTime int64) *T400 {
	return &T400{
		tlv:      NewTLV(0x0400, 0x0000, nil),
		key:      key,
		uin:      uin,
		guid:     guid,
		dpwd:     dpwd,
		appID:    appID,
		subAppID: subAppID,
		randSeed: randSeed,

		serverTime: serverTime,
	}
}

func (t *T400) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}

func (t *T400) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	if len(t.randSeed) == 0 {
		t.randSeed = make([]byte, 8)
	}
	v.WriteUint16(0x0001)
	v.WriteUint64(t.uin)
	v.WriteBytesL16V(t.guid[:])
	v.WriteBytesL16V(t.dpwd[:])
	v.WriteUint32(uint32(t.appID))
	v.WriteUint32(uint32(t.subAppID))
	v.WriteUint32(uint32(t.serverTime))
	v.WriteBytesL16V(t.randSeed)
	t.tlv.SetValue(bytes.NewBuffer(tea.NewCipher(t.key).Encrypt(v.Bytes())))
	return t.tlv.WriteTo(b)
}
