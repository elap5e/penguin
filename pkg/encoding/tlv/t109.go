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

type T109 struct {
	tlv  *TLV
	hash [16]byte
}

func NewT109(hash [16]byte) *T109 {
	return &T109{
		tlv:  NewTLV(0x0109, 0x0000, nil),
		hash: hash,
	}
}

func (t *T109) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	copy(t.hash[:], v.Bytes())
	return nil
}

func (t *T109) WriteTo(b *bytes.Buffer) error {
	t.tlv.SetValue(bytes.NewBuffer(t.hash[:]))
	return t.tlv.WriteTo(b)
}
