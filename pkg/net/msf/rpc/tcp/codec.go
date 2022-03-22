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

package tcp

import (
	"io"
	"net"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type codec struct {
	c    rpc.Client
	conn io.ReadWriteCloser

	buf *bytes.Buffer

	NetworkType uint8 // 0x00: Others; 0x01: Wi-Fi
	NetIPFamily uint8 // 0x00: Others; 0x01: IPv4; 0x02: IPv6; 0x03: Dual

	IMEI     string
	KSID     []byte
	IMSI     string
	Revision string
}

func NewCodec(conn net.Conn) rpc.Codec {
	return &codec{conn: conn}
}

func (c *codec) Close() error {
	return c.conn.Close()
}

var _ rpc.Codec = (*codec)(nil)
