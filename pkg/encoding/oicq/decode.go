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

package oicq

import (
	"fmt"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
)

func Unmarshal(p []byte, d *Data) error {
	n, err := checkValid(p)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(p[1 : n-1])
	if err := unmarshalHead(b, d); err != nil {
		return err
	}
	if err := unmarshalBody(b, d); err != nil {
		return err
	}
	return nil
}

func checkValid(v []byte) (int, error) {
	n := len(v)
	if v[0] != 0x02 {
		return 1, fmt.Errorf("unexpected prefix, got 0x%x", v[0])
	}
	if v[n-1] != 0x03 {
		return n, fmt.Errorf("unexpected suffix, got 0x%x", v[n-1])
	}
	return n, nil
}

func unmarshalHead(b *bytes.Buffer, d *Data) error {
	var err error
	if _, err = b.ReadUint16(); err != nil {
		return err
	}
	if d.Version, err = b.ReadUint16(); err != nil {
		return err
	}
	if d.ServiceMethod, err = b.ReadUint16(); err != nil {
		return err
	}
	if _, err = b.ReadUint16(); err != nil {
		return err
	}
	var uin uint32
	if uin, err = b.ReadUint32(); err != nil {
		return err
	}
	d.Uin = int64(uin)
	if _, err = b.ReadByte(); err != nil {
		return err
	}
	encryptMethod, err := b.ReadByte()
	if err != nil {
		return err
	}
	d.EncryptMethod = GetEncryptMethod(encryptMethod)
	if _, err = b.ReadByte(); err != nil {
		return err
	}
	return nil
}

func unmarshalBody(b *bytes.Buffer, d *Data) error {
	switch d.EncryptMethod {
	case EncryptMethod0x00:
	case EncryptMethod0x03:
		d.SharedSecret = d.RandomKey
	}
	tmp, err := tea.NewCipher(d.SharedSecret).Decrypt(b.Bytes())
	if err != nil {
		return err
	}
	b = bytes.NewBuffer(tmp)
	if err := unmarshalTLVs(b, d); err != nil {
		return err
	}
	return nil
}

func unmarshalTLVs(b *bytes.Buffer, d *Data) error {
	var err error
	if d.Type, err = b.ReadUint16(); err != nil {
		return err
	}
	if d.Code, err = b.ReadByte(); err != nil {
		return err
	}
	var l uint16
	if l, err = b.ReadUint16(); err != nil {
		return err
	}
	d.TLVs = map[uint16]tlv.Codec{}
	for i := 0; i < int(l); i++ {
		tlv := tlv.TLV{}
		tlv.ReadFrom(b)
		d.TLVs[tlv.GetType()] = &tlv
	}
	return nil
}
