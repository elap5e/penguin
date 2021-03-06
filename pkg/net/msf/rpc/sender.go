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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"strconv"
	"sync"
	"time"

	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Sender interface {
	Close() error
	Go(serviceMethod string, args *Args, reply *Reply, done chan *Call) *Call

	Run(ctx context.Context, cancel context.CancelFunc)
}

type sender struct {
	cancel context.CancelFunc

	c     Client
	codec Codec

	sendLock sync.Mutex
	sendChan chan struct{}
	sendList *list.List
	sendIdle bool

	mu       sync.Mutex
	pending  map[int32]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop

	heartbeating bool
	lastRecvTime time.Time
}

func NewSender(c Client, codec Codec) Sender {
	return &sender{
		c:        c,
		codec:    codec,
		sendLock: sync.Mutex{},
		sendChan: make(chan struct{}, 10),
		sendList: list.New(),
		pending:  make(map[int32]*Call),
	}
}

func (s *sender) recvLoop(ctx context.Context) {
	var err error
	var resp Response
	for err == nil {
		resp = Response{}
		err = s.codec.ReadResponseHeader(&resp)
		if err != nil {
			break
		}
		seq := resp.Seq
		s.mu.Lock()
		s.lastRecvTime = time.Now()
		call := s.pending[seq]
		delete(s.pending, seq)
		s.mu.Unlock()

		p, _ := json.Marshal(resp)
		log.Debug("recv:head:%s", p)
		info := fmt.Sprintf("ver:%d uin:%s seq:%d cmd:%s", resp.Version, resp.Username, resp.Seq, resp.ServiceMethod)
		if call != nil {
			err = s.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
			p, _ := json.Marshal(call.Reply)
			log.Trace("recv:%s data:\n%s", p, hex.Dump(call.Reply.Payload))
			log.Debug("recv:%s data:%d", p, len(call.Reply.Payload))
			log.Info("recv %s data:%d", info, len(call.Reply.Payload))
		} else {
			reply := &Reply{}
			err = s.codec.ReadResponseBody(reply)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
			p, _ := json.Marshal(reply)
			log.Trace("push:%s data:\n%s", p, hex.Dump(reply.Payload))
			log.Debug("push:%s data:%d", p, len(reply.Payload))
			log.Info("push %s data:%d", info, len(reply.Payload))
			go func() {
				args, err := s.c.Handle(reply.ServiceMethod, reply)
				if err == ErrCachedPush {
					log.Warn("cached push %s data:%d", info, len(reply.Payload))
				} else if err != nil {
					log.Error("handle %s error:%v", info, err)
					log.Debug("push:%s data:\n%s", p, hex.Dump(reply.Payload))
					log.Warn("skip %s data:%d", info, len(reply.Payload))
				} else if args != nil {
					s.goSend(args.ServiceMethod, args, nil, nil, true)
				}
			}()
		}
	}
	log.Error("recvLoop error: %v", err)
	s.loopError(err)
}

func (s *sender) sendLoop(ctx context.Context) {
	var err error
	var req Request
	for err == nil {
		s.sendLock.Lock()
		s.sendIdle = false
		elem := s.sendList.Front()
		if elem == nil {
			s.sendIdle = true
			s.sendLock.Unlock()
			<-s.sendChan // queue is empty, waiting for push unlock
			continue
		}
		call := elem.Value.(*Call)
		s.sendList.Remove(elem)
		s.sendLock.Unlock()

		req = Request{
			ServiceMethod: call.ServiceMethod,
			Seq:           call.Seq,
			Version:       call.Version,
			EncryptType:   EncryptTypeNotNeedEncrypt,
			Username:      strconv.FormatInt(call.Args.Uin, 10),
		}
		err := s.codec.WriteRequest(&req, call.Args)
		p, _ := json.Marshal(req)
		log.Debug("send:head:%s", p)
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
		p, _ = json.Marshal(call.Args)
		log.Trace("send:%s data:\n%s", p, hex.Dump(call.Args.Payload))
		log.Debug("send:%s data:%d", p, len(call.Args.Payload))
		log.Info("send ver:%d uin:%s seq:%d cmd:%s data:%d", req.Version, req.Username, req.Seq, req.ServiceMethod, len(call.Args.Payload))
	}
	log.Error("sendLoop error: %v", err)
	s.loopError(err)
}

func (s *sender) watchDog(ctx context.Context) {
	var now time.Time
	var err error
	duration := time.Second * 60
	ticker := time.NewTicker(duration)
	for err == nil {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
		}
		s.mu.Lock()
		now = time.Now()
		if s.lastRecvTime.Add(duration).After(now) {
			ticker.Reset(s.lastRecvTime.Add(duration).Sub(now))
			s.mu.Unlock()
			continue
		}
		s.mu.Unlock()
		if !s.heartbeating {
			if err := s.heartbeat(); err != nil {
				log.Error("watchDog sender, error: %v", err)
				s.Close()
			}
		}
		ticker.Reset(duration)
	}
}

func (s *sender) loopError(err error) {
	log.Error("loop error: %s", err)
	// Terminate pending calls.
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
	s.Close()
}

func (s *sender) Run(ctx context.Context, cancel context.CancelFunc) {
	s.cancel = func() { cancel() }
	go s.recvLoop(ctx)
	go s.sendLoop(ctx)
	go s.watchDog(ctx)
	select {
	case <-ctx.Done():
		return
	}
}

func (s *sender) push(call *Call, front bool) {
	// Register this call.
	s.mu.Lock()
	if s.shutdown || s.closing {
		s.mu.Unlock()
		call.Error = rpc.ErrShutdown
		call.done()
		return
	}
	if call.Reply != nil {
		s.pending[call.Seq] = call
	}
	s.mu.Unlock()

	s.sendLock.Lock()
	if front {
		s.sendList.PushFront(call)
	} else {
		s.sendList.PushBack(call)
	}
	if s.sendList.Len() != 0 && len(s.sendChan) == 0 && s.sendIdle {
		s.sendChan <- struct{}{}
	}
	s.sendLock.Unlock()
}

func (s *sender) Close() error {
	s.mu.Lock()
	defer s.mu.Lock()
	if s.closing {
		return fmt.Errorf("sender is closing")
	}
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
	if args.Seq == 0 {
		args.Seq = s.c.GetNextSeq()
	}
	call.Seq = args.Seq
	if args.Version == 0 {
		args.Version = VersionDefault
	}
	call.Version = args.Version
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

func (s *sender) heartbeat() (err error) {
	s.heartbeating = true
	ticker := time.NewTicker(time.Second)
	select {
	case <-ticker.C:
		ticker.Stop()
		err = ErrHeartbeatTimeout
	case call := <-s.goSend(service.MethodHeartbeatAlive, new(Args), new(Reply), make(chan *Call, 1), true).Done:
		err = call.Error
	}
	s.heartbeating = false
	return err
}

var _ Sender = (*sender)(nil)
