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

package rpc

import (
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/rpc"
	"strconv"
	"sync"
	"time"

	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Sender interface {
	Close() error
	Go(serviceMethod string, args *Args, reply *Reply, done chan *Call) *Call

	Run(ctx context.Context)
}

func NewSender(c Client, codec Codec) Sender {
	return &sender{
		c:       c,
		codec:   codec,
		send:    list.New(),
		wait:    list.New(),
		pending: make(map[int32]*Call),
	}
}

type sender struct {
	cancel context.CancelFunc

	c     Client
	codec Codec

	send   *list.List
	sendMu sync.Mutex
	wait   *list.List
	waitMu sync.Mutex

	beatAt    time.Time
	beating   bool
	beatRetry int32

	mu       sync.Mutex
	pending  map[int32]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}

func (s *sender) recvLoop(ctx context.Context) {
	var err error
	var resp Response
	for err == nil {
		select {
		case <-ctx.Done():
			return
		default:
		}

		resp = Response{}
		err = s.codec.ReadResponseHeader(&resp)
		if err != nil {
			return
		}
		seq := resp.Seq
		s.mu.Lock()
		call := s.pending[seq]
		delete(s.pending, seq)
		s.mu.Unlock()

		switch {
		case call == nil:
			// TODO: server push
			err = s.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		default:
			err = s.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
		p, _ := json.MarshalIndent(resp, "", "  ")
		log.Printf("sender.response:\n%s", string(p))
		p, _ = json.MarshalIndent(call.Reply, "", "  ")
		log.Printf("sender.call.Reply:\n%s", string(p))
		log.Printf("[R] uin:%s ver:%d seq:%d cmd:%s", resp.Username, resp.Version, resp.Seq, resp.ServiceMethod)
	}
	s.loopError(err)
}

func (s *sender) sendLoop(ctx context.Context) {
	var err error
	var req Request
	for err == nil {
		select {
		case <-ctx.Done():
			return
		default:
		}

		s.sendMu.Lock()
		elem := s.send.Front()
		if elem == nil {
			s.sendMu.Unlock()
			continue
		}
		call := elem.Value.(*Call)
		s.send.Remove(elem)
		s.sendMu.Unlock()

		req = Request{
			ServiceMethod: call.ServiceMethod,
			Seq:           call.Seq,
			Version:       VersionDefault,
			EncryptType:   EncryptTypeNotNeedEncrypt,
			Username:      strconv.FormatInt(call.Args.Uin, 10),
		}
		err := s.codec.WriteRequest(&req, call.Args)
		if err != nil {
			s.mu.Lock()
			call = s.pending[call.Seq]
			delete(s.pending, call.Seq)
			s.mu.Unlock()
			if call != nil {
				call.Error = err
				call.done()
			}
		}
		p, _ := json.MarshalIndent(req, "", "  ")
		log.Printf("sender.requset:\n%s", string(p))
		p, _ = json.MarshalIndent(call.Args, "", "  ")
		log.Printf("sender.call.Args:\n%s", string(p))
		log.Printf("[S] uin:%s ver:%d seq:%d cmd:%s", req.Username, req.Version, req.Seq, req.ServiceMethod)
	}
	s.loopError(err)
}

func (s *sender) waitLoop(ctx context.Context) {
	// wait for more calls to arrive
}

func (s *sender) loopError(err error) {
	log.Println("loop error:", err)
	// Terminate pending calls.
	s.sendMu.Lock()
	s.mu.Lock()
	s.shutdown = true
	closing := s.closing
	if err == io.EOF {
		if closing {
			err = rpc.ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range s.pending {
		call.Error = err
		call.done()
	}
	s.mu.Unlock()
	s.sendMu.Unlock()
	s.cancel()
}

func (s *sender) Run(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)
	go s.recvLoop(ctx)
	go s.sendLoop(ctx)
	go s.waitLoop(ctx)
	select {
	case <-ctx.Done():
		return
	}
}

func (s *sender) push(call *Call, front bool) {
	s.sendMu.Lock()
	defer s.sendMu.Unlock()

	// Register this call.
	s.mu.Lock()
	if s.shutdown || s.closing {
		s.mu.Unlock()
		call.Error = rpc.ErrShutdown
		call.done()
		return
	}
	s.pending[call.Seq] = call
	s.mu.Unlock()

	if front {
		s.send.PushFront(call)
	} else {
		s.send.PushBack(call)
	}
}

func (s *sender) Close() error {
	s.cancel()
	return s.codec.Close()
}

func (s *sender) Go(serviceMethod string, args *Args, reply *Reply, done chan *Call) *Call {
	return s.goSend(serviceMethod, args, reply, done, false)
}

func (s *sender) goSend(serviceMethod string, args *Args, reply *Reply, done chan *Call, front bool) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	args.ServiceMethod = serviceMethod
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
	s.push(call, front)
	return call
}

func (s *sender) heartbeat() {
	s.beatAt = time.Now()
	s.beating = true
	s.beatRetry++
	args := &Args{
		Uin:     0,
		Seq:     s.c.GetNextSeq(),
		Payload: []byte{0x00, 0x00, 0x00, 0x04},
	}
	call := <-s.goSend(service.MethodHeartbeatAlive, args, new(Reply), make(chan *Call, 1), true).Done
	s.beating = false
	if call.Error != nil {
		return
	}
}

var _ Sender = (*sender)(nil)
