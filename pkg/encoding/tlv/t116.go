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

type T116 struct {
	*TLV
	miscBitmap uint32
	subSigMap  uint32
	subAppIDs  []uint64
}

func NewT116(miscBitmap, subSigMap uint32, subAppIDs []uint64) *T116 {
	return &T116{
		TLV:        NewTLV(0x0116, 0x0000, nil),
		miscBitmap: miscBitmap,
		subSigMap:  subSigMap,
		subAppIDs:  subAppIDs,
	}
}

func (t *T116) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.ReadByte(); err != nil {
		return err
	}
	if t.miscBitmap, err = v.ReadUint32(); err != nil {
		return err
	}
	if t.subSigMap, err = v.ReadUint32(); err != nil {
		return err
	}
	l, err := v.ReadByte()
	if err != nil {
		return err
	}
	t.subAppIDs = make([]uint64, l)
	for i := range t.subAppIDs {
		j, err := v.ReadUint32()
		if err != nil {
			return err
		}
		t.subAppIDs[i] = uint64(j)
	}
	return nil
}

func (t *T116) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteByte(0)
	v.WriteUint32(t.miscBitmap)
	v.WriteUint32(t.subSigMap)
	v.WriteByte(uint8(len(t.subAppIDs)))
	for i := range t.subAppIDs {
		v.WriteUint32(uint32(t.subAppIDs[i]))
	}
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
