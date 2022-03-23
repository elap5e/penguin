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

type T124 struct {
	*TLV
	osType []byte
	osVer  []byte
	nwType uint16
	simOp  []byte
	bArr4  []byte
	apn    []byte
}

func NewT124(osType, osVer []byte, nwType uint16, simOp, bArr4, apn []byte) *T124 {
	return &T124{
		TLV:    NewTLV(0x0124, 0x0000, nil),
		osType: osType,
		osVer:  osVer,
		nwType: nwType,
		simOp:  simOp,
		bArr4:  bArr4,
		apn:    apn,
	}
}

func (t *T124) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	_, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}

func (t *T124) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteBytesL16V(t.osType, 0x0010)
	v.WriteBytesL16V(t.osVer, 0x0010)
	v.WriteUint16(t.nwType)
	v.WriteBytesL16V(t.simOp, 0x0010)
	v.WriteBytesL16V(t.bArr4, 0x0020)
	v.WriteBytesL16V(t.apn, 0x0010)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
