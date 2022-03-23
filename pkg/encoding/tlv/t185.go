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

type T185 struct {
	tlv *TLV
	i   uint8
}

func NewT185(i uint8) *T185 {
	return &T185{
		tlv: NewTLV(0x0185, 0x0000, nil),
		i:   i,
	}
}

func (t *T185) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	ver, err := v.ReadByte()
	if err != nil {
		return err
	} else if ver != 0x01 {
		return fmt.Errorf("version 0x%x not support", ver)
	}
	if t.i, err = v.ReadByte(); err != nil {
		return err
	}
	return nil
}

func (t *T185) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteByte(0x01)
	v.WriteByte(t.i)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
