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

package highway

import (
	"encoding/binary"
	"errors"
	"io"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/constant"
	"github.com/elap5e/penguin/pkg/bytes"
	"github.com/elap5e/penguin/pkg/net/highway/pb"
)

type Codec struct {
	conn io.ReadWriteCloser

	bufHead *bytes.Buffer
	payload *bytes.Buffer
	reply   *Reply
}

func NewCodec(conn io.ReadWriteCloser) *Codec {
	return &Codec{
		conn:    conn,
		bufHead: bytes.NewBuffer([]byte{}),
		payload: bytes.NewBuffer([]byte{}),
		reply:   new(Reply),
	}
}

func (c *Codec) Close() error {
	return c.conn.Close()
}

func (c *Codec) read() (int64, error) {
	flag := make([]byte, 1)
	if _, err := c.conn.Read(flag); err != nil {
		return 0, err
	}
	if flag[0] != 0x28 {
		return 1, errors.New("invalid flag")
	}

	c.bufHead.Reset()
	if n, err := io.CopyN(c.bufHead, c.conn, 4); err != nil {
		return n + 1, err
	}
	l1 := int64(binary.BigEndian.Uint32(c.bufHead.Bytes()))
	c.bufHead.Reset()
	c.bufHead.Grow(int(l1))

	c.payload.Reset()
	if n, err := io.CopyN(c.payload, c.conn, 4); err != nil {
		return n + 5, err
	}
	l2 := int64(binary.BigEndian.Uint32(c.payload.Bytes()))
	c.payload.Reset()
	c.payload.Grow(int(l2))

	if n, err := io.CopyN(c.bufHead, c.conn, l1); err != nil {
		return n + 9, err
	}
	if n, err := io.CopyN(c.payload, c.conn, l2); err != nil {
		return n + l1 + 9, err
	}

	if _, err := c.conn.Read(flag); err != nil {
		return l1 + l2 + 9, err
	}
	if flag[0] != 0x29 {
		return l1 + l2 + 10, errors.New("invalid flag")
	}
	return l1 + l2 + 10, nil
}

func (c *Codec) ReadResponseHeader(resp *Response) (err error) {
	if _, err = c.read(); err != nil {
		return err
	}
	head := pb.CSDataHighwayHead_RspDataHighwayHead{}
	if err = proto.Unmarshal(c.bufHead.Bytes(), &head); err != nil {
		return err
	}
	resp.ServiceMethod = string(head.GetMsgBasehead().GetCommand())
	c.reply.ServiceMethod = resp.ServiceMethod
	resp.Seq = int32(head.GetMsgBasehead().GetSeq())
	c.reply.Seq = resp.Seq
	resp.Username = string(head.GetMsgBasehead().GetUin())
	c.reply.Uin, _ = strconv.ParseInt(resp.Username, 10, 64)
	resp.Code = int32(head.GetErrorCode())
	c.reply.Code = resp.Code
	c.reply.Payload = c.payload.Bytes()
	return nil
}

func (c *Codec) ReadResponseBody(reply *Reply) (err error) {
	reply.Uin = c.reply.Uin
	reply.Seq = c.reply.Seq
	reply.Code = c.reply.Code
	reply.ServiceMethod = c.reply.ServiceMethod
	reply.Payload = append([]byte{}, c.reply.Payload...)
	return nil
}

func (c *Codec) WriteRequest(req *Request, args *Args) error {
	head, err := proto.Marshal(&pb.CSDataHighwayHead_ReqDataHighwayHead{
		MsgBasehead: &pb.CSDataHighwayHead_DataHighwayHead{
			Version:    1,
			Uin:        []byte(req.Username),
			Command:    []byte(req.ServiceMethod),
			Seq:        uint32(req.Seq),
			RetryTimes: 0,
			Appid:      uint32(req.AppID),
			Dataflag:   0x00001000,
			CommandId:  uint32(req.CommandID),
			BuildVer:   []byte(""),
			LocaleId:   constant.LocaleID,
			EnvId:      0,
		},
		MsgSeghead:    args.SegHead,
		ReqExtendinfo: []byte{},
		Timestamp:     0,
	})
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(0x28)
	buf.WriteUint32(uint32(len(head)))
	buf.WriteUint32(uint32(len(args.Payload)))
	buf.Write(head)
	buf.Write(args.Payload)
	buf.WriteByte(0x29)
	if _, err := buf.WriteTo(c.conn); err != nil {
		return err
	}
	return nil
}
