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

type T107 struct {
	*TLV
	i  uint16
	i2 uint8
	i3 uint16
	i4 uint8
}

func NewT107(i uint16, i2 uint8, i3 uint16, i4 uint8) *T107 {
	return &T107{
		TLV: NewTLV(0x0107, 0x0000, nil),
		i:   i,
		i2:  i2,
		i3:  i3,
		i4:  i4,
	}
}

func (t *T107) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	if t.i, err = v.ReadUint16(); err != nil {
		return err
	}
	if t.i2, err = v.ReadByte(); err != nil {
		return err
	}
	if t.i3, err = v.ReadUint16(); err != nil {
		return err
	}
	if t.i4, err = v.ReadByte(); err != nil {
		return err
	}
	return nil
}

func (t *T107) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(t.i)
	v.WriteByte(t.i2)
	v.WriteUint16(t.i3)
	v.WriteByte(t.i4)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
