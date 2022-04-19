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
	"strconv"
	"strings"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/log"
)

type T511 struct {
	*TLV
	domains []string
}

func NewT511(domains []string) *T511 {
	return &T511{
		TLV:     NewTLV(0x0511, 0x0000, nil),
		domains: domains,
	}
}

func (t *T511) ReadFrom(b *bytes.Buffer) error {
	if err := t.TLV.ReadFrom(b); err != nil {
		return err
	}
	_, err := t.TLV.GetValue()
	if err != nil {
		return err
	}
	panic("not implement")
}

func (t *T511) WriteTo(b *bytes.Buffer) error {
	v := bytes.NewBuffer([]byte{})
	var domains []string
	for i := range t.domains {
		if t.domains[i] != "" {
			domains = append(domains, t.domains[i])
		}
	}
	v.WriteUint16(uint16(len(domains)))
	var flag uint8
	for _, domain := range domains {
		idx0 := strings.Index(domain, "(")
		idx1 := strings.Index(domain, ")")
		if idx0 != 0 || idx1 <= 0 {
			flag = 0x01
		} else {
			i, err := strconv.Atoi(domain[idx0+1 : idx1])
			if err != nil {
				log.Printf("GetT511 error: %s", err.Error())
			}
			var z1 = (1048576 & i) > 0
			var z2 = (i & 134217728) > 0
			if z1 {
				flag = 0x01
			} else {
				flag = 0x00
			}
			if z2 {
				flag |= 0x02
			}
			domain = domain[idx1+1:]
		}
		v.WriteByte(flag)
		v.WriteStringL16V(domain)
	}
	t.TLV.SetValue(v)
	return t.TLV.WriteTo(b)
}
