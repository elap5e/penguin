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
	"net"

	"github.com/elap5e/penguin/pkg/bytes"
)

type T1 struct {
	*TLV
	uin int64
	ip  net.IP

	serverTime int64
}

func NewT1(uin int64, ip net.IP, serverTime int64) *T1 {
	return &T1{
		TLV: NewTLV(0x0001, 0x0014, nil),
		uin: uin,
		ip:  ip,

		serverTime: serverTime,
	}
}

func (t *T1) GetUin() (int64, error) {
	return t.uin, nil
}

func (t *T1) GetIP() (net.IP, error) {
	return t.ip, nil
}

func (t *T1) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	ipVer, err := v.ReadUint16()
	if err != nil {
		return err
	} else if ipVer != 0x0001 {
		return fmt.Errorf("ip version 0x%x not support", ipVer)
	}
	if _, err := v.ReadUint32(); err != nil {
		return err
	}
	uin, err := v.ReadUint32()
	if err != nil {
		return err
	}
	t.uin = int64(uin)
	if _, err := v.ReadUint32(); err != nil {
		return err
	}
	t.ip = make([]byte, 4)
	if _, err := v.Read(t.ip); err != nil {
		return err
	}
	if _, err := v.ReadUint32(); err != nil {
		return err
	}
	return nil
}

func (t *T1) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	v.WriteInt16(0x0001)
	v.WriteInt32(random.Int31())
	v.WriteInt32(int32(t.uin))
	v.WriteInt32(int32(t.serverTime))
	v.Write(t.ip.To4())
	v.WriteInt16(0x0000)
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
