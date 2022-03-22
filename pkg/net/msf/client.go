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
	"math/rand"
	"sync/atomic"

	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func NewClient() rpc.Client {
	return &client{
		rs:  rpc.NewSender(),
		seq: rand.Int31n(100000),
	}
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

func (c *client) GetAppID() int32 {
	return -1
}

func (c *client) SetAppID(id int32) {
}

func (c *client) GetFakeApp(username string) *rpc.FakeApp {
	return nil
}

func (c *client) GetTickets(username string) rpc.Tickets {
	return nil
}

var _ rpc.Client = (*client)(nil)
