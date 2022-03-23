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

type T100 struct {
	*TLV
	appID    uint64
	subAppID uint64
	acVer    uint32
	sigMap   uint32

	ssoVer uint32
}

func NewT100(appID, subAppID uint64, acVer, sigMap, ssoVer uint32) *T100 {
	return &T100{
		TLV:      NewTLV(0x0100, 0x0000, nil),
		appID:    appID,
		subAppID: subAppID,
		acVer:    acVer,
		sigMap:   sigMap,

		ssoVer: ssoVer,
	}
}

func (t *T100) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.ReadUint16(); err != nil {
		return err
	}
	if _, err = v.ReadUint32(); err != nil {
		return err
	}
	appID, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	subAppID, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.subAppID = uint64(subAppID)
	if t.acVer, err = v.ReadUint32(); err != nil {
		return err
	}
	if t.sigMap, err = v.ReadUint32(); err != nil {
		return err
	}
	return nil
}

func (t *T100) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint16(0x0001)
	v.WriteUint32(t.ssoVer)
	v.WriteUint32(uint32(t.appID))
	v.WriteUint32(uint32(t.subAppID))
	v.WriteUint32(t.acVer)
	v.WriteUint32(t.sigMap)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
