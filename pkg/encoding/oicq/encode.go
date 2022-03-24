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
	"encoding/binary"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
)

func Marshal(d *Data) ([]byte, error) {
	h, err := marshalHead(d)
	if err != nil {
		return nil, err
	}
	b, err := marshalBody(d)
	if err != nil {
		return nil, err
	}
	p := append(h, b...)
	binary.BigEndian.PutUint16(p[1:], uint16(len(p)))
	return p, nil
}

func marshalHead(d *Data) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	b.WriteByte(0x02)
	b.WriteUint16(0x0000)
	b.WriteUint16(d.Version)
	b.WriteUint16(d.ServiceMethod)
	b.WriteUint16(0x0001)
	b.WriteUint32(uint32(d.Uin))
	b.WriteByte(0x03)
	switch d.EncryptMethod {
	case EncryptMethodECDH:
		b.WriteByte(0x07)
	case EncryptMethodECDH0x87:
		b.WriteByte(0x87)
	case EncryptMethodST:
		b.WriteByte(0x45)
	}
	b.WriteByte(0x00)
	b.WriteUint32(0x00000002)
	b.WriteUint32(0x00000000)
	b.WriteUint32(0x00000000)
	return b.Bytes(), nil
}

func marshalBody(d *Data) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	switch d.EncryptMethod {
	case EncryptMethodECDH, EncryptMethodECDH0x87:
		b.WriteByte(0x02)
		b.WriteByte(0x01)
		b.Write(d.RandomKey[:])
		b.WriteUint16(0x0131)
		b.WriteUint16(uint16(d.KeyVersion))
		b.WriteUint16(uint16(len(d.PublicKey)))
		b.Write(d.PublicKey)
	case EncryptMethodST:
		b.WriteByte(0x01)
		b.WriteByte(0x03)
		b.Write(d.RandomKey[:])
		b.WriteUint16(0x0102)
		b.WriteUint16(0x0000)
		d.SharedSecret = d.RandomKey
	}
	data, err := marshalTLVs(d)
	if err != nil {
		return nil, err
	}
	b.Write(tea.NewCipher(d.SharedSecret).Encrypt(data))
	b.WriteByte(0x03)
	return b.Bytes(), nil
}

func marshalTLVs(d *Data) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	b.WriteUint16(d.Type)
	for i := range d.TLVs {
		if d.TLVs[i] == nil {
			delete(d.TLVs, i)
		}
	}
	b.WriteUint16(uint16(len(d.TLVs)))
	for i := range d.TLVs {
		d.TLVs[i].WriteTo(b)
	}
	return b.Bytes(), nil
}
