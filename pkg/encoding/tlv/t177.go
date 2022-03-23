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

type T177 struct {
	tlv        *TLV
	buildTime  int64
	sdkVersion string
}

func NewT177(buildTime int64, sdkVersion string) *T177 {
	return &T177{
		tlv:        NewTLV(0x0177, 0x0000, nil),
		buildTime:  buildTime,
		sdkVersion: sdkVersion,
	}
}

func (t *T177) ReadFrom(b *bytes.Buffer) error {
	if err := t.tlv.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.tlv.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.ReadByte(); err != nil {
		return err
	}
	buildTime, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.buildTime = int64(buildTime)
	if t.sdkVersion, err = v.ReadStringL16V(); err != nil {
		return err
	}
	return nil
}

func (t *T177) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteByte(0x01)
	v.WriteUint32(uint32(t.buildTime))
	v.WriteStringL16V(t.sdkVersion)
	t.tlv.SetValue(v)
	return t.tlv.WriteTo(b)
}
