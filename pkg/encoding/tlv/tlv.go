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
	"math/rand"
	"time"

	"github.com/elap5e/penguin/pkg/bytes"
)

var random = rand.New(rand.NewSource(time.Now().UTC().UnixMicro()))

type Codec interface {
	ReadFrom(b *bytes.Buffer) error
	WriteTo(b *bytes.Buffer) error
}

type TLV struct {
	t uint16
	l uint16
	v *bytes.Buffer
}

func (v *TLV) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"t%x(%d)\"", v.t, v.l)), nil
}

func NewTLV(t uint16, l uint16, v *bytes.Buffer) *TLV {
	return &TLV{t: t, l: l, v: v}
}

func (t *TLV) SetValue(v *bytes.Buffer) {
	t.v = v
}

func (t *TLV) GetType() uint16 {
	return t.t
}

func (t *TLV) GetValue() (*bytes.Buffer, error) {
	return t.v, nil
}

func (t *TLV) MustGetValue() *bytes.Buffer {
	return t.v
}

func (t *TLV) ReadFrom(b *bytes.Buffer) error {
	var err error
	if t.t, err = b.ReadUint16(); err != nil {
		return err
	}
	if t.l, err = b.ReadUint16(); err != nil {
		return err
	}
	v := make([]byte, t.l)
	if _, err = b.Read(v); err != nil {
		return err
	}
	t.v = bytes.NewBuffer(v)
	return nil
}

func (t *TLV) WriteTo(b *bytes.Buffer) error {
	v := t.v.Bytes()
	t.l = uint16(len(v))
	b.WriteUint16(t.t)
	b.WriteUint16(t.l)
	b.Write(v)
	return nil
}
