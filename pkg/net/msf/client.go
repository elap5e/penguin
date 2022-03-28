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
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elap5e/penguin/pkg/encoding/tlv"
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
	r := rand.New(rand.NewSource(uin))
	buf := make([]byte, 20)
	_, err := r.Read(buf)
	if err != nil {
		log.Fatalf("failed to generate device config")
	}
	ipv4 := net.IPv4(192, 168, 0, buf[0])
	mac1 := fmt.Sprintf("00:50:%02x:%02x:00:%02x", buf[1], buf[2], buf[0])
	mac2 := fmt.Sprintf("00:50:%02x:%02x:00:%02x", buf[1], buf[2], buf[3])
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", buf[4:7], buf[8:9], buf[10:11], buf[12:13], buf[14:19])
	imei := fmt.Sprintf("86030802%07d", r.Int31n(10000000))
	imsi := fmt.Sprintf("460001%09d", r.Int31n(1000000000))
	osid := fmt.Sprintf("RKQ1.%07d.002", r.Int31n(10000000))
	return &rpc.FakeSource{
		App: &rpc.FakeApp{
			FixID:      537114575,
			AppID:      537114575,
			PkgName:    "com.tencent.mobileqq",
			VerName:    "8.8.83", // 8.8.83.7540
			Revision:   "8.8.83.7b3989f8",
			SigHash:    [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d},
			BuildAt:    1645432578,
			SDKVer:     "6.0.0.2497",
			SSOVer:     18,
			ImageType:  1,
			MiscBitMap: 0b00001000111101111111111101111100, // 150470524
			CanCaptcha: true,
		},
		Device: &rpc.FakeDevice{
			OS: rpc.FakeDeviceOS{
				Type:        "android",
				Version:     "11",
				BuildBrand:  []byte("Xiaomi"),
				BuildModel:  "Redmi K20",
				BuildID:     osid,
				SDKVersion:  30,
				NetworkType: 2,
			},
			APNName:       []byte("wifi"),
			SIMOPName:     []byte("CMCC"),
			Bootloader:    "unknown",
			ProcVersion:   "Linux version 2.6.18-92.el5 (brewbuilder@ls20-bc2-13.build.redhat.com)",
			Codename:      "davinci",
			Incremental:   "20.10.20",
			Fingerprint:   "Xiaomi/davinci/davinci:11/" + osid + "/20.10.20:user/release-keys",
			BootID:        uuid,
			Baseband:      "4.3.c5-00069-SM6150_GEN_PACK-1",
			InnerVersion:  "20.10.20",
			NetworkType:   1,
			NetIPFamily:   3,
			IPv4Address:   ipv4,
			IPv6Address:   ipv4,
			MACAddress:    mac1,
			BSSIDAddress:  mac2,
			SSIDAddress:   "unknown",
			IMEI:          imei,
			IMSI:          imsi,
			GUID:          md5.Sum(append([]byte(osid), mac1...)),
			GUIDFlag:      uint32((1 << 24 & 0xff000000) | (0 << 8 & 0xff00)),
			IsGUIDFileNil: false,
			IsGUIDGenSucc: true,
			IsGUIDChanged: false,
		},
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
