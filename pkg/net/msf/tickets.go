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

package msf

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func newTickets(uin int64) *rpc.Tickets {
	var tickets rpc.Tickets
	tickets = rpc.Tickets{
		A1:       &rpc.Ticket{},
		A2:       &rpc.Ticket{},
		A5:       &rpc.Ticket{},
		A8:       &rpc.Ticket{},
		D2:       &rpc.Ticket{},
		LSKey:    &rpc.Ticket{},
		SKey:     &rpc.Ticket{},
		SID:      &rpc.Ticket{},
		Sig64:    &rpc.Ticket{},
		SuperKey: &rpc.Ticket{},
		ST:       &rpc.Ticket{},
		STWeb:    &rpc.Ticket{},
		VKey:     &rpc.Ticket{},
		Domains:  map[string]string{},
		KSID:     []byte{},
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Read(tickets.A1.Key[:])
	return &tickets
}

func getTickets(uin int64) *rpc.Tickets {
	file := ".penguin/tickets/" + strconv.FormatInt(uin, 10) + ".json"
	data, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		data, err = json.MarshalIndent(newTickets(uin), "", "  ")
		if err == nil {
			err = ioutil.WriteFile(file, data, 0644)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}
	var tickets rpc.Tickets
	err = json.Unmarshal(data, &tickets)
	if err != nil {
		log.Fatalln(err)
	}
	return &tickets
}

func setTickets(uin int64, tickets *rpc.Tickets, tlvs map[uint16]tlv.Codec) {
	iss := time.Now().Unix()
	for k, v := range tlvs {
		v := v.(*tlv.TLV)
		switch k {
		// case 0x011d:
		// case 0x011f:
		// case 0x0130:
		// case 0x0133:
		// case 0x0134:
		// case 0x0163:
		// case 0x016a:
		// case 0x0522:
		// case 0x0528:
		// case 0x0537:
		// case 0x0550:
		case 0x0102:
			tickets.A8.Sig = v.MustGetValue().Bytes()
		case 0x0103:
			tickets.STWeb.Sig = v.MustGetValue().Bytes()
		case 0x0106:
			tickets.A1.Sig = v.MustGetValue().Bytes()
		case 0x0108:
			tickets.KSID = v.MustGetValue().Bytes()
		case 0x010a:
			tickets.A2.Sig = v.MustGetValue().Bytes()
		case 0x010b:
			tickets.A5.Sig = v.MustGetValue().Bytes()
			tickets.A5.Iss = iss
			tickets.A5.Exp = -1
		case 0x010c:
			copy(tickets.A1.Key[:], v.MustGetValue().Bytes())
		case 0x010d:
			copy(tickets.A2.Key[:], v.MustGetValue().Bytes())
		case 0x010e:
			copy(tickets.ST.Key[:], v.MustGetValue().Bytes())
		case 0x0114:
			tickets.ST.Sig = v.MustGetValue().Bytes()
			tickets.ST.Iss = iss
			tickets.ST.Exp = -1
		case 0x0118:
			log.Trace("t%x main_display_name:\n%s", k, hex.Dump(v.MustGetValue().Bytes()))
		case 0x011a:
			log.Trace("t%x face, age, gender, nick:\n%s", k, hex.Dump(v.MustGetValue().Bytes()))
		case 0x011c:
			tickets.LSKey.Sig = v.MustGetValue().Bytes()
		case 0x0120:
			tickets.SKey.Sig = v.MustGetValue().Bytes()
		case 0x0121:
			tickets.Sig64.Sig = v.MustGetValue().Bytes()
			tickets.Sig64.Iss = iss
			tickets.Sig64.Exp = -1
		case 0x0136:
			tickets.VKey.Sig = v.MustGetValue().Bytes()
		case 0x0138:
			buf := bytes.NewBuffer(v.MustGetValue().Bytes())
			l, _ := buf.ReadUint32()
			for i := 0; i < int(l); i++ {
				key, _ := buf.ReadUint16()
				exp, _ := buf.ReadUint32()
				_, _ = buf.ReadUint32()
				switch key {
				case 0x0102:
					tickets.A8.Iss = iss
					tickets.A8.Exp = iss + int64(exp)
				case 0x0103:
					tickets.STWeb.Iss = iss
					tickets.STWeb.Exp = iss + int64(exp)
				case 0x0106:
					tickets.A1.Iss = iss
					tickets.A1.Exp = iss + int64(exp)
				case 0x010a:
					tickets.A2.Iss = iss
					tickets.A2.Exp = iss + int64(exp)
				case 0x011c:
					tickets.LSKey.Iss = iss
					tickets.LSKey.Exp = iss + int64(exp)
				case 0x0120:
					tickets.SKey.Iss = iss
					tickets.SKey.Exp = iss + int64(exp)
				case 0x0136:
					tickets.VKey.Iss = iss
					tickets.VKey.Exp = iss + int64(exp)
				case 0x0143:
					tickets.D2.Iss = iss
					tickets.D2.Exp = iss + int64(exp)
				case 0x0164:
					tickets.SID.Iss = iss
					tickets.SID.Exp = iss + int64(exp)
				default:
					log.Trace("t%x change time not parsed:%d", key, exp)
				}
			}
		case 0x0143:
			tickets.D2.Sig = v.MustGetValue().Bytes()
		case 0x0164:
			tickets.SID.Sig = v.MustGetValue().Bytes()
		case 0x016d:
			tickets.SuperKey.Sig = v.MustGetValue().Bytes()
			tickets.SuperKey.Iss = iss
			tickets.SuperKey.Exp = -1
		case 0x0305:
			copy(tickets.D2.Key[:], v.MustGetValue().Bytes())
		case 0x0512:
			if tickets.Domains == nil {
				tickets.Domains = make(map[string]string)
			}
			buf := bytes.NewBuffer(v.MustGetValue().Bytes())
			l, _ := buf.ReadUint16()
			for i := 0; i < int(l); i++ {
				key, _ := buf.ReadStringL16V()
				tickets.Domains[key], _ = buf.ReadStringL16V()
				_, _ = buf.ReadUint16()
			}
		default:
			log.Trace("t%x not parsed:\n%s", k, hex.Dump(v.MustGetValue().Bytes()))
		}
	}
	file := ".penguin/tickets/" + strconv.FormatInt(uin, 10) + ".json"
	data, err := json.MarshalIndent(tickets, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(file, data, 0644)
	}
}
