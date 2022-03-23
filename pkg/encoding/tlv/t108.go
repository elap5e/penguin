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

type T108 struct {
	*TLV
	ksid []byte
}

func NewT108(ksid []byte) *T108 {
	return &T108{
		TLV:  NewTLV(0x0108, 0x0000, nil),
		ksid: ksid,
	}
}

func (t *T108) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	t.ksid = v.Bytes()
	return nil
}

func (t *T108) WriteTo(b *bytes.Buffer) error {
	t.TLV.SetValue(bytes.NewBuffer(t.ksid))
	return t.TLV.WriteTo(b)
}
