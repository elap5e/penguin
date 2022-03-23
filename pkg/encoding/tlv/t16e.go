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

type T16E struct {
	tlv   *TLV
	model []byte
}

func NewT16E(model []byte) *T16E {
	return &T16E{
		tlv:   NewTLV(0x016e, 0x0000, nil),
		model: model,
	}
}

func (t *T16E) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.model = v.Bytes()
	return nil
}

func (t *T16E) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteBytesL16V(t.model, 0x0040)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
