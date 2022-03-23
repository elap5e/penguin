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

type T201 struct {
	tlv   *TLV
	bArr  []byte
	bArr2 []byte
	bArr3 []byte
	bArr4 []byte
}

func NewT201(bArr, bArr2, bArr3, bArr4 []byte) *T201 {
	return &T201{
		tlv:   NewTLV(0x0201, 0x0000, nil),
		bArr:  bArr,
		bArr2: bArr2,
		bArr3: bArr3,
		bArr4: bArr4,
	}
}

func (t *T201) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if t.bArr, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	if t.bArr2, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	if t.bArr3, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	if t.bArr4, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T201) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteBytesL16V(t.bArr)
	v.WriteBytesL16V(t.bArr2)
	v.WriteBytesL16V(t.bArr3)
	v.WriteBytesL16V(t.bArr4)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}