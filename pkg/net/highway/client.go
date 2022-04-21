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
	"crypto/md5"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/highway/pb"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type Client struct {
	codec *Codec

	reqMutex sync.Mutex // protects following
	req      Request

	mu       sync.Mutex // protects following
	seq      int32
	pending  map[int32]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}

func Dial(network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), nil
}

func NewClient(conn io.ReadWriteCloser) *Client {
	client := &Client{
		codec:   NewCodec(conn),
		seq:     random.Int31n(100000),
		pending: make(map[int32]*Call),
	}
	go client.input()
	return client
}

func (c *Client) getNextSeq() int32 {
	seq := atomic.AddInt32(&c.seq, 1)
	if seq > 1000000 {
		c.seq = random.Int31n(100000) + 60000
	}
	return seq
}

func (c *Client) send(call *Call) {
	c.reqMutex.Lock()
	defer c.reqMutex.Unlock()

	c.mu.Lock()
	if c.shutdown || c.closing {
		c.mu.Unlock()
		call.Error = rpc.ErrShutdown
		call.done()
		return
	}
	seq := c.seq
	c.seq++
	c.pending[seq] = call
	c.mu.Unlock()

	c.req.Username = strconv.FormatInt(call.Args.Uin, 10)
	c.req.Seq = seq
	c.req.ServiceMethod = call.ServiceMethod
	c.req.CommandID = call.Args.CommandID
	c.req.AppID = 0 // TODO: fix this
	err := c.codec.WriteRequest(&c.req, call.Args)
	if err != nil {
		c.mu.Lock()
		call = c.pending[seq]
		delete(c.pending, seq)
		c.mu.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (c *Client) input() {
	var err error
	var resp Response
	for err == nil {
		resp = Response{}
		err = c.codec.ReadResponseHeader(&resp)
		if err != nil {
			break
		}
		seq := resp.Seq
		c.mu.Lock()
		call := c.pending[seq]
		delete(c.pending, seq)
		c.mu.Unlock()

		if call != nil {
			err = c.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		} else {
			err = c.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		}
	}

	c.reqMutex.Lock()
	c.mu.Lock()
	c.shutdown = true
	closing := c.closing
	if err == io.EOF {
		if closing {
			err = rpc.ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
	c.mu.Unlock()
	c.reqMutex.Unlock()
}

func (c *Client) Close() error {
	c.mu.Lock()
	if c.closing {
		c.mu.Unlock()
		return rpc.ErrShutdown
	}
	c.closing = true
	c.mu.Unlock()
	return c.codec.Close()
}

func (c *Client) Go(serviceMethod string, args *Args, reply *Reply, done chan *Call) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	args.ServiceMethod = serviceMethod
	if args.Seq == 0 {
		args.Seq = c.getNextSeq()
	}
	call.Seq = args.Seq
	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *Call, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	c.send(call)
	return call
}

func (c *Client) Call(serviceMethod string, args *Args, reply *Reply) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

func (c *Client) UploadFile(uin int64, cmd int32, name string, ticket []byte) error {
	f, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.UploadReadSeeker(uin, cmd, f, ticket)
}

func (c *Client) UploadReadSeeker(uin int64, cmd int32, r io.ReadSeeker, ticket []byte) error {
	hash := md5.New()
	size, err := io.Copy(hash, r)
	if err != nil {
		return err
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return c.UploadReader(uin, cmd, size, r, hash.Sum(nil), ticket)
}

func (c *Client) UploadReader(uin int64, cmd int32, size int64, r io.Reader, digest, ticket []byte) error {
	if err := c.echo(uin); err != nil {
		return err
	}
	chunk, offset := make([]byte, 0x00080000), int64(0) // 512KiB
	for {
		n, err := r.Read(chunk)
		if err == io.EOF || (err == nil && n == 0) {
			break
		}
		if err != nil {
			return err
		}
		if err := c.uploadChunk(uin, cmd, size, offset, chunk[:n], digest, ticket); err != nil {
			return err
		}
		offset += int64(n)
	}
	return nil
}

func (c *Client) echo(uin int64) error {
	return c.Call(ServiceMethodEcho, &Args{Uin: uin}, new(Reply))
}

func (c *Client) uploadChunk(uin int64, cmd int32, size, offset int64, chunk, digest, ticket []byte) error {
	hash := md5.New()
	hash.Write(chunk)
	segHead := pb.CSDataHighwayHead_SegHead{
		Serviceid:     0, // nil
		Filesize:      uint64(size),
		Dataoffset:    uint64(offset),
		Datalength:    uint32(len(chunk)),
		Rtcode:        0, // nil
		Serviceticket: ticket,
		Flag:          0, // nil
		Md5:           hash.Sum(nil),
		FileMd5:       digest,
		CacheAddr:     0, // nil
		QueryTimes:    0, // nil
		UpdateCacheip: 0, // nil
		CachePort:     0, // nil
	}
	return c.Call(ServiceMethodUpload, &Args{Uin: uin, CommandID: cmd, SegHead: &segHead, Payload: chunk}, new(Reply))
}
