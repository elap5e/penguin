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
	"fmt"

	"github.com/elap5e/penguin/pkg/bytes"
)

type T141 struct {
	tlv     *TLV
	simOp   []byte
	netType uint16
	apn     []byte
}

func NewT141(simOp []byte, netType uint16, apn []byte) *T141 {
	return &T141{
		tlv:     NewTLV(0x0141, 0x0000, nil),
		simOp:   simOp,
		netType: netType,
		apn:     apn,
	}
}

func (t *T141) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	ver, err := v.ReadUint16()
	if err != nil {
		return err
	} else if ver != 0x0001 {
		return fmt.Errorf("version 0x%x not support", ver)
	}
	if t.simOp, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	if t.netType, err = v.ReadUint16(); err != nil {
		return err
	}
	if t.apn, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T141) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(0x0001)
	v.WriteBytesL16V(t.simOp)
	v.WriteUint16(t.netType)
	v.WriteBytesL16V(t.apn)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
