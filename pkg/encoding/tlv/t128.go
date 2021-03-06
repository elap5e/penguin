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

type T128 struct {
	*TLV
	isGuidFileNil bool
	isGuidGenSucc bool
	isGuidChanged bool
	guidFlag      uint32
	model         []byte
	guid          []byte
	brand         []byte
}

func NewT128(isGuidFileNil, isGuidGenSucc, isGuidChanged bool, guidFlag uint32, model, guid, brand []byte) *T128 {
	return &T128{
		TLV:           NewTLV(0x0128, 0x0000, nil),
		isGuidFileNil: isGuidFileNil,
		isGuidGenSucc: isGuidGenSucc,
		isGuidChanged: isGuidChanged,
		guidFlag:      guidFlag,
		model:         model,
		guid:          guid,
		brand:         brand,
	}
}

func (t *T128) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.ReadUint16(); err != nil {
		return err
	}
	if t.isGuidFileNil, err = v.ReadBool(); err != nil {
		return err
	}
	if t.isGuidGenSucc, err = v.ReadBool(); err != nil {
		return err
	}
	if t.isGuidChanged, err = v.ReadBool(); err != nil {
		return err
	}
	if t.guidFlag, err = v.ReadUint32(); err != nil {
		return err
	}
	if t.model, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	if t.guid, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	if t.brand, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T128) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(0x0000)
	v.WriteBool(t.isGuidFileNil)
	v.WriteBool(t.isGuidGenSucc)
	v.WriteBool(t.isGuidChanged)
	v.WriteUint32(t.guidFlag)
	v.WriteBytesL16V(t.model, 0x0020)
	v.WriteBytesL16V(t.guid, 0x0010)
	v.WriteBytesL16V(t.brand, 0x0010)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
