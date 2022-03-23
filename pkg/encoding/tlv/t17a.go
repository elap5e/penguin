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

type T17A struct {
	*TLV
	appID uint64
}

func NewT17A(appID uint64) *T17A {
	return &T17A{
		TLV:   NewTLV(0x017a, 0x0000, nil),
		appID: appID,
	}
}

func (t *T17A) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	appID, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.appID = uint64(appID)
	return nil
}

func (t *T17A) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteUint32(uint32(t.appID))
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
