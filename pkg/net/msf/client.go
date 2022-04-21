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
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/penguin/fake"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/rpc/tcp"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type Client struct {
	rpc.Sender
	seq int32

	mu          sync.Mutex
	handlers    map[string]rpc.Handler
	sessions    map[int64]*rpc.Session
	stickets    map[int64]*rpc.Tickets
	fakeSources map[int64]*fake.Source
}

func NewClient(ctx context.Context) rpc.Client {
	c := &Client{
		seq:         random.Int31n(100000),
		mu:          sync.Mutex{},
		handlers:    make(map[string]rpc.Handler),
		sessions:    make(map[int64]*rpc.Session),
		stickets:    make(map[int64]*rpc.Tickets),
		fakeSources: make(map[int64]*fake.Source),
	}
	if err := c.DialDefault(ctx); err != nil {
		log.Warn("msf.Client.DialDefault error: %v", err)
	}
	return c
}

func (c *Client) Dial(network, address string) error {
	return c.DialWithContext(context.Background(), network, address)
}

func (c *Client) DialDefault(ctx context.Context) error {
	return c.DialWithContext(ctx, "tcp", "msfwifi.3g.qq.com:8080")
}

func (c *Client) DialWithContext(ctx context.Context, network, address string) error {
	if c.Sender != nil {
		if err := c.Sender.Close(); err != nil {
			return err
		}
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	c.Sender = rpc.NewSender(c, tcp.NewCodec(c, conn))
	go c.Sender.Run(ctx)
	return nil
}

func (c *Client) Close() error {
	return c.Sender.Close()
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
	c.mu.Lock()
	key := strings.ToLower(serviceMethod)
	if _, ok := c.handlers[key]; ok {
		c.mu.Unlock()
		return fmt.Errorf("service method %s already registered", serviceMethod)
	}
	c.handlers[key] = handler
	c.mu.Unlock()
	return nil
}

func (c *Client) GetNextSeq() int32 {
	seq := atomic.AddInt32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = random.Int31n(100000) + 60000
	}
	return seq
}

func (c *Client) GetFakeSource(uin int64) *fake.Source {
	c.mu.Lock()
	fakeSource := c.fakeSources[uin]
	if fakeSource == nil {
		c.fakeSources[uin] = fake.NewMobileQQAndroidSource(uin)
		fakeSource = c.fakeSources[uin]
	}
	c.mu.Unlock()
	return fakeSource
}

func (c *Client) GetServerTime() int64 {
	return time.Now().Unix()
}

func (c *Client) GetSession(uin int64) *rpc.Session {
	c.mu.Lock()
	session := c.sessions[uin]
	if session == nil {
		c.sessions[uin] = getSession(uin)
		session = c.sessions[uin]
	}
	c.mu.Unlock()
	return session
}

func (c *Client) SetSession(uin int64, tlvs map[uint16]tlv.Codec) {}

func (c *Client) SetSessionAuth(uin int64, auth []byte) {
	session := c.GetSession(uin)
	c.mu.Lock()
	session.Auth = auth
	c.mu.Unlock()
}

func (c *Client) SetSessionCookie(uin int64, cookie []byte) {
	session := c.GetSession(uin)
	c.mu.Lock()
	session.Cookie = append([]byte{}, cookie...)
	c.mu.Unlock()
}

func (c *Client) GetTickets(uin int64) *rpc.Tickets {
	c.mu.Lock()
	tickets := c.stickets[uin]
	if tickets == nil {
		c.stickets[uin] = getTickets(uin)
		tickets = c.stickets[uin]
	}
	c.mu.Unlock()
	return tickets
}

func (c *Client) SetTickets(uin int64, tlvs map[uint16]tlv.Codec) {
	c.mu.Lock()
	tickets := c.stickets[uin]
	if tickets == nil {
		c.stickets[uin] = getTickets(uin)
		tickets = c.stickets[uin]
	}
	setTickets(uin, tickets, tlvs)
	c.mu.Unlock()
}

var _ rpc.Client = (*Client)(nil)
