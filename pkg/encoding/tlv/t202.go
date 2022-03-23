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

type T202 struct {
	tlv      *TLV
	md5BSSID [16]byte
	ssid     []byte
}

func NewT202(md5BSSID [16]byte, ssid []byte) *T202 {
	return &T202{
		tlv:      NewTLV(0x0202, 0x0000, nil),
		md5BSSID: md5BSSID,
		ssid:     ssid,
	}
}

func (t *T202) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	md5BSSID, err := v.ReadBytesL16V()
	if err != nil {
		return err
	}
	copy(t.md5BSSID[:], md5BSSID)
	if t.ssid, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T202) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteBytesL16V(t.md5BSSID[:], 0x0010)
	v.WriteBytesL16V(t.ssid, 0x0020)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
