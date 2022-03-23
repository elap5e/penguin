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
	"github.com/elap5e/penguin/pkg/encoding/tlv"
)

type Data struct {
	Version       uint16
	ServiceMethod uint16
	Uin           int64
	EncryptMethod EncryptMethod
	RandomKey     [16]byte
	KeyVersion    int16
	PublicKey     []byte
	SharedSecret  [16]byte
	Type          uint16
	Code          uint8
	TLVs          map[uint16]tlv.Codec
}

type EncryptMethod uint8

var (
	EncryptMethod0x00 EncryptMethod = 0x00
	EncryptMethod0x03 EncryptMethod = 0x03
	EncryptMethodECDH EncryptMethod = 0x07 | 0x80 // 0x07: no password login?
	EncryptMethodST   EncryptMethod = 0x45
	EncryptMethodNULL EncryptMethod = 0xff
)

func GetEncryptMethod(v uint8) EncryptMethod {
	switch v {
	case 0x00:
		return EncryptMethod0x00
	case 0x03:
		return EncryptMethod0x03
	case 0x07, 0x87:
		return EncryptMethodECDH
	case 0x45:
		return EncryptMethodST
	default:
		return EncryptMethodNULL
	}
}
