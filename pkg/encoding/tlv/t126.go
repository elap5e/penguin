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

type T126 struct {
	*TLV
	random []byte
}

func NewT126(random []byte) *T126 {
	return &T126{
		TLV:    NewTLV(0x0126, 0, nil),
		random: random,
	}
}

func (t *T126) ReadFrom(b *bytes.Buffer) error {
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
	if t.random, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T126) GetRandom() ([]byte, error) {
	return t.random, nil
}

func (t *T126) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(0)
	v.WriteBytesL16V(t.random)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
