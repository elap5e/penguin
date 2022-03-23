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

type T10A struct {
	tlv *TLV
	a2  []byte // A2
}

func NewT10A(a2 []byte) *T10A {
	return &T10A{
		tlv: NewTLV(0x010a, 0x0000, nil),
		a2:  a2,
	}
}

func (t *T10A) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	t.a2 = v.Bytes()
	return nil
}

func (t *T10A) WriteTo(b *bytes.Buffer) error {
	t.tlv.SetValue(bytes.NewBuffer(t.a2))
	return t.tlv.WriteTo(b)
}
