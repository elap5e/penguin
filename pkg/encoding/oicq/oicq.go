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
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Data struct {
	Version       uint16               `json:"version,omitempty"`
	ServiceMethod uint16               `json:"service_method,omitempty"`
	Uin           int64                `json:"uin,omitempty"`
	EncryptMethod EncryptMethod        `json:"encrypt_method,omitempty"`
	RandomKey     rpc.Key16Bytes       `json:"random_key,omitempty"`
	KeyVersion    int16                `json:"key_version,omitempty"`
	PublicKey     []byte               `json:"public_key,omitempty"`
	SharedSecret  rpc.Key16Bytes       `json:"shared_secret,omitempty"`
	Type          uint16               `json:"type,omitempty"`
	Code          uint8                `json:"code,omitempty"`
	TLVs          map[uint16]tlv.Codec `json:"tlvs,omitempty"`
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
