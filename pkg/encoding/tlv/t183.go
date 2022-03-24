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

type T183 struct {
	*TLV
	salt uint64
}

func NewT183(salt uint64) *T183 {
	return &T183{
		TLV:  NewTLV(0x0183, 0x0000, nil),
		salt: salt,
	}
}

func (t *T183) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	if t.salt, err = v.ReadUint64(); err != nil {
		return err
	}
	return nil
}

func (t *T183) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint64(t.salt)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
