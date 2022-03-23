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

type T18 struct {
	tlv   *TLV
	appID uint64
	acVer uint32
	uin   int64
	i2    uint16
}

func NewT18(appID uint64, acVer uint32, uin int64, i2 uint16) *T18 {
	return &T18{
		tlv:   NewTLV(0x0018, 0x0000, nil),
		appID: appID,
		acVer: acVer,
		uin:   uin,
		i2:    i2,
	}
}

func (t *T18) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	pingVer, err := v.ReadUint16()
	if err != nil {
		return err
	} else if pingVer != 0x0001 {
		return fmt.Errorf("ping version 0x%x not support", pingVer)
	}
	ssoVer, err := v.ReadUint32()
	if err != nil {
		return err
	} else if ssoVer != 0x00000600 {
		return fmt.Errorf("sso version 0x%x not support", ssoVer)
	}
	appID, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	if t.acVer, err = v.ReadUint32(); err != nil {
		return err
	}
	uin, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.uin = int64(uin)
	if t.i2, err = v.ReadUint16(); err != nil {
		return err
	}
	if _, err = v.ReadUint16(); err != nil {
		return err
	}
	return nil
}

func (t *T18) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(0x0001)
	v.WriteUint32(0x00000600)
	v.WriteUint32(uint32(t.appID))
	v.WriteUint32(t.acVer)
	v.WriteUint32(uint32(t.uin))
	v.WriteUint16(t.i2)
	v.WriteUint16(0x0000)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
