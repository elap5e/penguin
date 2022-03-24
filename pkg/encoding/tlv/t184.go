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
	"crypto/md5"
	"encoding/binary"

	"github.com/elap5e/penguin/pkg/bytes"
)

type T184 struct {
	*TLV
	salt     uint64
	password string
}

func NewT184(salt uint64, password string) *T184 {
	return &T184{
		TLV:      NewTLV(0x0184, 0x0000, nil),
		salt:     salt,
		password: password,
	}
}

func (t *T184) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	_, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}

func (t *T184) WriteTo(b *bytes.Buffer) error {
	v := md5.Sum([]byte(t.password))
	tmp := append(v[:], make([]byte, 8)...)
	binary.BigEndian.PutUint64(tmp[16:], t.salt)
	v = md5.Sum(tmp)
	t.TLV.SetValue(bytes.NewBuffer(v[:]))
	return t.TLV.WriteTo(b)
}
