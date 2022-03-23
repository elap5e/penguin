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

type T8 struct {
	*TLV
	i1       uint16
	localeID uint32
	i3       uint16
}

func NewT8(i1 uint16, localeID uint32, i3 uint16) *T8 {
	return &T8{
		TLV:      NewTLV(0x0008, 0x0000, nil),
		i1:       i1,
		localeID: localeID,
		i3:       i3,
	}
}

func (t *T8) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	if t.i1, err = v.ReadUint16(); err != nil {
		return err
	}
	if t.localeID, err = v.ReadUint32(); err != nil {
		return err
	}
	if t.i3, err = v.ReadUint16(); err != nil {
		return err
	}
	return nil
}

func (t *T8) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(t.i1)
	v.WriteUint32(t.localeID)
	v.WriteUint16(t.i3)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
