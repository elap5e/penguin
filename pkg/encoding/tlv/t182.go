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

type T182 struct {
	*TLV
	msgCnt    uint16
	timeLimit uint16
}

func NewT182(msgCnt, timeLimit uint16) *T182 {
	return &T182{
		TLV:       NewTLV(0x0182, 0x0000, nil),
		msgCnt:    msgCnt,
		timeLimit: timeLimit,
	}
}

func (t *T182) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	if _, err = v.ReadByte(); err != nil {
		return err
	}
	if t.msgCnt, err = v.ReadUint16(); err != nil {
		return err
	}
	if t.timeLimit, err = v.ReadUint16(); err != nil {
		return err
	}
	return nil
}

func (t *T182) GetMessageCount() (uint16, error) {
	return t.msgCnt, nil
}

func (t *T182) GetTimeLimit() (uint16, error) {
	return t.timeLimit, nil
}

func (t *T182) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteByte(0)
	v.WriteUint16(t.msgCnt)
	v.WriteUint16(t.timeLimit)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
