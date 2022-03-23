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

type T142 struct {
	tlv   *TLV
	apkID []byte
}

func NewT142(apkID []byte) *T142 {
	return &T142{
		tlv:   NewTLV(0x0142, 0x0000, nil),
		apkID: apkID,
	}
}

func (t *T142) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if ver, err := v.ReadUint16(); err != nil {
		return err
	} else if ver != 0x0001 {
		return fmt.Errorf("version 0x%x not support", ver)
	}
	if t.apkID, err = v.ReadBytesL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T142) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(0x0000)
	v.WriteBytesL16V(t.apkID, 0x0020)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
