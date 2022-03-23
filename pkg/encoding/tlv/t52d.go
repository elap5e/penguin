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
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/encoding/tlv/pb"
)

type T52D struct {
	*TLV
	deviceReport *pb.DeviceReport
}

func NewT52D(deviceReport *pb.DeviceReport) *T52D {
	return &T52D{
		TLV:          NewTLV(0x052d, 0x0000, nil),
		deviceReport: deviceReport,
	}
}

func (t *T52D) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	v, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	deviceInfo := pb.DeviceReport{}
	if err := proto.Unmarshal(v.Bytes(), &deviceInfo); err != nil {
		return err
	}
	panic("not implement")
}

func (t *T52D) WriteTo(b *bytes.Buffer) error {
	v, _ := proto.Marshal(t.deviceReport)
	t.TLV.SetValue(bytes.NewBuffer(v))
	return t.TLV.WriteTo(b)
}
