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

type T2 struct {
	tlv  *TLV
	code []byte
	sign []byte
}

func NewT2(code, sign []byte) *T2 {
	return &T2{
		tlv:  NewTLV(0x0002, 0x0000, nil),
		code: code,
		sign: sign,
	}
}

func (t *T2) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	sigVer, err := v.ReadUint16()
	if err != nil {
		return err
	} else if sigVer != 0x0000 {
		return fmt.Errorf("sig version 0x%x not support", sigVer)
	}
	if t.code, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	if t.sign, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T2) GetCode() ([]byte, error) {
	return t.code, nil
}

func (t *T2) GetSign() ([]byte, error) {
	return t.sign, nil
}

func (t *T2) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(0x0000)
	v.WriteBytesL16V(t.code)
	v.WriteBytesL16V(t.sign)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
