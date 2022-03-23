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
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/crypto/tea"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (c *codec) read() (int64, error) {
	// Read the first 4 bytes to get the length of the response.
	c.buf.Reset()
	if n, err := io.CopyN(c.buf, c.conn, 4); err != nil {
		return n, err
	}
	l := int64(binary.BigEndian.Uint32(c.buf.Bytes())) - 4
	// Read the next l bytes to get the response.
	c.buf.Grow(int(l))
	n, err := io.CopyN(c.buf, c.conn, l)
	return n + 4, err
}

func (c *codec) ReadResponseHeader(resp *rpc.Response) (err error) {
	if _, err = c.read(); err != nil {
		return err
	}

	log.Printf("dump of read:\n%s", hex.Dump(c.buf.Bytes()))
	// Skip the first 4 bytes for the length of the response.
	if _, err = c.buf.ReadUint32(); err != nil {
		return err
	}
	// Read next 4 bytes to get the version.
	if resp.Version, err = c.buf.ReadUint32(); err != nil {
		return err
	}
	if resp.Version != rpc.VersionDefault && resp.Version != rpc.VersionSimple {
		return fmt.Errorf("tcp: unsupported version 0x%x", resp.Version)
	}
	// Read next byte to get the encrpyt type.
	if resp.EncryptType, err = c.buf.ReadByte(); err != nil {
		return err
	}
	if resp.EncryptType != rpc.EncryptTypeNotNeedEncrypt &&
		resp.EncryptType != rpc.EncryptTypeEncryptByD2Key &&
		resp.EncryptType != rpc.EncryptTypeEncryptByZeros {
		return fmt.Errorf("tcp: unsupported encrypt type 0x%x", resp.EncryptType)
	}
	// Skip the next unused byte.
	if _, err = c.buf.ReadByte(); err != nil {
		return err
	}
	// Read next 4 bytes to get the username length and read the username.
	if resp.Username, err = c.buf.ReadStringL32(); err != nil {
		return err
	}
	c.reply.Uin, _ = strconv.ParseInt(resp.Username, 10, 64)
	// Decrypt the response body if the response body is encrypted.
	switch resp.EncryptType {
	case rpc.EncryptTypeNotNeedEncrypt:
		// c.buf = bytes.NewBuffer(c.buf.Bytes())
	case rpc.EncryptTypeEncryptByD2Key:
		buf, err := tea.NewCipher(c.cl.GetTickets(c.reply.Uin).D2.Key).Decrypt(c.buf.Bytes())
		if err != nil {
			return err
		}
		c.buf = bytes.NewBuffer(buf)
	case rpc.EncryptTypeEncryptByZeros:
		buf, err := tea.NewCipher([16]byte{}).Decrypt(c.buf.Bytes())
		if err != nil {
			return err
		}
		c.buf = bytes.NewBuffer(buf)
	}

	var n uint32
	log.Printf("dump of read body:\n%s", hex.Dump(c.buf.Bytes()))
	// Read the first 4 bytes for the length of the response body header.
	if n, err = c.buf.ReadUint32(); err != nil {
		return err
	}
	// Calculate the length of the response body payload.
	n = uint32(c.buf.Len()) - n + 4
	// Read the next 4 bytes for the sequence.
	if c.reply.Seq, err = c.buf.ReadInt32(); err != nil {
		return err
	}
	resp.Seq = c.reply.Seq
	// Read the next 4 bytes for the status code.
	if c.reply.Code, err = c.buf.ReadInt32(); err != nil {
		return err
	}
	// Read the next 4 bytes for the status message length and read the status message.
	if c.reply.Message, err = c.buf.ReadStringL32(); err != nil {
		return err
	}
	// Read the next 4 bytes for the service method length and read the service method.
	if c.reply.ServiceMethod, err = c.buf.ReadStringL32(); err != nil {
		return err
	}
	resp.ServiceMethod = c.reply.ServiceMethod
	// Read the next 4 bytes for the cookie length and read the cookie.
	if c.reply.Cookie, err = c.buf.ReadBytesL32(); err != nil {
		return err
	}
	// Read the next 4 bytes for the flag.
	if c.reply.Flag, err = c.buf.ReadUint32(); err != nil {
		return err
	}
	if c.reply.Flag != rpc.FlagNoCompression && c.reply.Flag != rpc.FlagZlibCompression {
		return fmt.Errorf("tcp: unsupported flag 0x%x", c.reply.Flag)
	}
	// Read iff the buffer is larger than the length of the response body payload.
	if c.buf.Len() > int(n) {
		if c.reply.ReserveField, err = c.buf.ReadBytesL32(); err != nil {
			return err
		}
	}
	// Read the next 4 bytes for the length of the payload and read the payload.
	if c.reply.Payload, err = c.buf.ReadBytesL32(); err != nil {
		return err
	}
	// Decompress the payload if the payload is compressed.
	switch c.reply.Flag {
	case rpc.FlagNoCompression:
		// c.reply.Payload = c.reply.Payload
	case rpc.FlagZlibCompression:
		reader, err := zlib.NewReader(bytes.NewBuffer(c.reply.Payload))
		if err != nil {
			return err
		}
		defer reader.Close()
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, reader); err != nil {
			return err
		}
		c.reply.Payload = buf.Bytes()
	}
	return nil
}

func (c *codec) ReadResponseBody(reply *rpc.Reply) (err error) {
	reply.Uin = c.reply.Uin
	reply.Seq = c.reply.Seq
	reply.Code = c.reply.Code
	reply.Message = c.reply.Message
	reply.ServiceMethod = c.reply.ServiceMethod
	reply.Cookie = c.reply.Cookie
	reply.Flag = c.reply.Flag
	reply.ReserveField = c.reply.ReserveField
	reply.Payload = c.reply.Payload
	return nil
}
