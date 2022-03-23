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
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"

	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/rpc/tcp"
)

func NewClient(ctx context.Context) rpc.Client {
	conn, _ := net.Dial("tcp", "14.22.5.202:8080")
	c := &client{seq: rand.Int31n(100000)}
	c.rs = rpc.NewSender(c, tcp.NewCodec(c, conn))
	go c.rs.Run(ctx)
	return c
}

type client struct {
	rs  rpc.Sender
	seq int32
}

func (c *client) Close() error {
	return c.rs.Close()
}

func (c *client) Go(serviceMethod string, args *rpc.Args, reply *rpc.Reply, done chan *rpc.Call) *rpc.Call {
	return c.rs.Go(serviceMethod, args, reply, done)
}

func (c *client) Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done
	return call.Error
}

func (c *client) GetNextSeq() int32 {
	seq := atomic.AddInt32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = rand.Int31n(100000) + 60000
	}
	return seq
}

func (c *client) GetFake(uin int64) *rpc.Fake {
	// TODO: hardcoded for now
	return &rpc.Fake{
		FixID:    537044845,
		AppID:    537044845,
		Revision: "8.8.83.7b3989f8",
		Device: &rpc.FakeDevice{
			OS: rpc.FakeDeviceOS{
				Type:        "android",
				Version:     "11",
				BuildBrand:  []byte("Xiaomi"),
				BuildModel:  "Redmi K20",
				SDKVersion:  uint32(30),
				NetworkType: 0x0002,
			},
			APNName:     []byte("wifi"),
			SIMOPName:   []byte("CMCC"),
			NetworkType: 0x01,
			NetIPFamily: 0x03,
			IMEI:        fmt.Sprintf("86030802%07d", rand.Int31n(10000000)),
			IMSI:        fmt.Sprintf("460001%09d", rand.Int31n(1000000000)),
		},
	}
}

func (c *client) GetServerTime() int64 {
	return 0
}

func (c *client) GetSession(uin int64) *rpc.Session {
	return &rpc.Session{
		Auth:         []byte{},
		Cookie:       []byte{},
		KSID:         []byte{},
		RandomKey:    [16]byte{},
		KeyVersion:   0,
		PublicKey:    []byte{},
		SharedSecret: [16]byte{},
	}
}

func (c *client) GetTickets(uin int64) *rpc.Tickets {
	return &rpc.Tickets{
		A1: rpc.Ticket{
			Key: [16]byte{},
			Sig: []byte{},
		},
		A2: rpc.Ticket{
			Key: [16]byte{},
			Sig: []byte{},
		},
		D2: rpc.Ticket{
			Key: [16]byte{},
			Sig: []byte{},
		},
	}
}

func (c *client) SetSession(uin int64, tlvs map[uint16]tlv.Codec) {}
func (c *client) SetSessionAuth(uin int64, auth []byte)           {}
func (c *client) SetSessionKSID(uin int64, ksid []byte)           {}
func (c *client) SetTickets(uin int64, tlvs map[uint16]tlv.Codec) {}

var _ rpc.Client = (*client)(nil)
