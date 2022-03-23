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
)

type T147 struct {
	tlv             *TLV
	appID           uint64
	apkVersion      []byte
	md5APKSignature [16]byte
}

func NewT147(appID uint64, apkVersion []byte, md5APKSignature [16]byte) *T147 {
	return &T147{
		tlv:             NewTLV(0x0147, 0x0000, nil),
		appID:           appID,
		apkVersion:      apkVersion,
		md5APKSignature: md5APKSignature,
	}
}

func (t *T147) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	appID, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	if t.apkVersion, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	md5APKSignature, err := v.ReadBytesL16V()
	if err != nil {
		return err
	}
	copy(t.md5APKSignature[:], md5APKSignature)
	return nil
}

func (t *T147) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint32(uint32(t.appID))
	v.WriteBytesL16V(t.apkVersion, 0x0020)
	v.WriteBytesL16V(t.md5APKSignature[:], 0x0020)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
