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

type T544 struct {
	*TLV
	uin    int64
	guid   [16]byte
	sdkVer string
	typ    uint16
}

func NewT544(uin int64, guid [16]byte, sdkVer string, typ uint16) *T544 {
	return &T544{
		TLV:    NewTLV(0x0544, 0x0000, nil),
		uin:    uin,
		guid:   guid,
		sdkVer: sdkVer,
		typ:    typ,
	}
}

func (t *T544) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	panic("not implement")
}

func (t *T544) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint32(uint32(t.uin))
	v.WriteBytesL16V(t.guid[:])
	v.WriteStringL16V(t.sdkVer)
	v.WriteUint32(uint32(t.typ))
	v.WriteUint32(0)
	panic("not implement")
}
