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
	"strings"
	"sync/atomic"
	"time"

	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/net/msf/fake"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/rpc/tcp"
)

func NewClient(ctx context.Context) rpc.Client {
	conn, _ := net.Dial("tcp", "msfwifi.3g.qq.com:8080")
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	cl := &Client{
		seq:      rr.Int31n(100000),
		handlers: make(map[string]rpc.Handler),
		sessions: make(map[int64]*rpc.Session),
		stickets: make(map[int64]*rpc.Tickets),
	}
	cl.rs = rpc.NewSender(cl, tcp.NewCodec(cl, conn))
	go cl.rs.Run(ctx)
	return cl
}

type Client struct {
	rs  rpc.Sender
	seq int32

	handlers map[string]rpc.Handler
	sessions map[int64]*rpc.Session
	stickets map[int64]*rpc.Tickets
}

func (c *Client) Close() error {
	return c.rs.Close()
}

func (c *Client) Go(serviceMethod string, args *rpc.Args, reply *rpc.Reply, done chan *rpc.Call) *rpc.Call {
	return c.rs.Go(serviceMethod, args, reply, done)
}

func (c *Client) Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done
	return call.Error
}

func (c *Client) Handle(serviceMethod string, reply *rpc.Reply) (*rpc.Args, error) {
	if handler, ok := c.handlers[strings.ToLower(serviceMethod)]; ok {
		return handler(reply)
	}
	return nil, rpc.ErrNotHandled
}

func (c *Client) Register(serviceMethod string, handler rpc.Handler) error {
	key := strings.ToLower(serviceMethod)
	if _, ok := c.handlers[key]; ok {
		return fmt.Errorf("service method %s already registered", serviceMethod)
	}
	c.handlers[key] = handler
	return nil
}

func (c *Client) GetNextSeq() int32 {
	seq := atomic.AddInt32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = rand.Int31n(100000) + 60000
	}
	return seq
}

func (c *Client) GetFakeSource(uin int64) *rpc.FakeSource {
	return &rpc.FakeSource{
		App:    fake.AndroidMobileQQApp,
		Device: fake.NewAndroidDevice(uin),
	}
}

func (c *Client) GetServerTime() int64 {
	return time.Now().Unix()
}

func (c *Client) GetSession(uin int64) *rpc.Session {
	session := c.sessions[uin]
	if session == nil {
		c.sessions[uin] = getSession(uin)
		session = c.sessions[uin]
	}
	return session
}

func (c *Client) SetSession(uin int64, tlvs map[uint16]tlv.Codec) {}

func (c *Client) SetSessionAuth(uin int64, auth []byte) {
	c.GetSession(uin).Auth = auth
}

func (c *Client) SetSessionCookie(uin int64, cookie []byte) {
	c.GetSession(uin).Cookie = cookie
}

func (c *Client) GetTickets(uin int64) *rpc.Tickets {
	tickets := c.stickets[uin]
	if tickets == nil {
		c.stickets[uin] = getTickets(uin)
		tickets = c.stickets[uin]
	}
	return tickets
}

func (c *Client) SetTickets(uin int64, tlvs map[uint16]tlv.Codec) {
	tickets := c.stickets[uin]
	if tickets == nil {
		c.stickets[uin] = getTickets(uin)
		tickets = c.stickets[uin]
	}
	setTickets(uin, tickets, tlvs)
}

var _ rpc.Client = (*Client)(nil)
