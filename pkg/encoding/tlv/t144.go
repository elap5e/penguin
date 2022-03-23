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

type T144 struct {
	tlv  *TLV
	key  [16]byte // Key
	tlvs []Codec
}

func NewT144(key [16]byte, tlvs ...Codec) *T144 {
	return &T144{
		tlv:  NewTLV(0x0144, 0x0000, nil),
		key:  key,
		tlvs: tlvs,
	}
}

func (t *T144) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	_, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}

func (t *T144) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(uint16(len(t.tlvs)))
	for i := range t.tlvs {
		t.tlvs[i].WriteTo(v)
	}
	t.tlv.SetValue(bytes.NewBuffer(tea.NewCipher(t.key).Encrypt(v.Bytes())))
	return t.tlv.WriteTo(b)
}
