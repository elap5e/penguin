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

package daemon

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/config"
	"github.com/elap5e/penguin/daemon/account"
	"github.com/elap5e/penguin/daemon/auth"
	"github.com/elap5e/penguin/daemon/channel"
	"github.com/elap5e/penguin/daemon/chat"
	"github.com/elap5e/penguin/daemon/contact"
	"github.com/elap5e/penguin/daemon/message"
	"github.com/elap5e/penguin/daemon/service"
	"github.com/elap5e/penguin/fake"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

var random = rand.New(rand.NewSource(time.Now().UTC().UnixMicro()))

type Daemon struct {
	ctx context.Context
	cfg *config.Config
	uin int64

	c rpc.Client

	accm *account.Manager
	athm *auth.Manager
	chtm *chat.Manager
	chnm *channel.Manager
	cntm *contact.Manager
	msgm *message.Manager
	svcm *service.Manager

	mu      sync.Mutex
	isReady bool
	readyCh chan struct{}
	msgChan chan *penguin.Message

	wg      sync.WaitGroup
	waiting bool

	heartbeating bool
	lastRecvTime time.Time
}

func New(ctx context.Context, cfg *config.Config) *Daemon {
	d := &Daemon{
		ctx: ctx,
		cfg: cfg,
		c:   msf.NewClient(),
		mu:  sync.Mutex{},
	}
	d.accm = account.NewManager(d.ctx, d, account.NewMemStore())
	d.athm = auth.NewManager(d.ctx, d.c, d)
	d.chtm = chat.NewManager(d.ctx, d, chat.NewMemStore())
	d.chnm = channel.NewManager(d.ctx, d, channel.NewMemStore())
	d.cntm = contact.NewManager(d.ctx, d, contact.NewMemStore())
	d.msgm = message.NewManager(d.ctx, d)
	d.svcm = service.NewManager(d.ctx, d)
	d.msgChan = make(chan *penguin.Message, 1000)
	d.readyCh = make(chan struct{}, 10)
	d.Waiting()
	return d
}

func (d *Daemon) watchDog(uin int64) {
	var now time.Time
	var err error
	duration := time.Second * 270
	ticker := time.NewTicker(duration)
	for err == nil {
		<-ticker.C
		d.mu.Lock()
		now = time.Now()
		if d.lastRecvTime.Add(duration).After(now) {
			ticker.Reset(d.lastRecvTime.Add(duration).Sub(now))
			d.mu.Unlock()
			continue
		}
		d.mu.Unlock()
		if !d.heartbeating {
			if _, err = d.svcm.RegisterAppRegister(uin); err != nil {
				log.Error("watchDog daemon, error: %v", err)
			}
		}
		ticker.Reset(duration)
	}
}

func (d *Daemon) Run() error {
	log.Info("rpc client is starting")
	errCh := make(chan error, 1)
	go func() {
		errCh <- d.run()
	}()
	<-d.readyCh
	log.Info("rpc client is ready, start the daemon")
	if err := d.init(); err != nil {
		return err
	}
	return <-errCh
}

func (d *Daemon) run() error {
	var err error
	for {
		ctx, cancel := context.WithCancel(context.Background())
		d.c.BindCancelFunc(cancel)
		if err = d.c.DialDefault(ctx); err != nil {
			log.Warn("daemon.Daemon.run error: %v", err)
			break
		}
		if d.checkReady() {
			if err = d.resync(d.uin); err != nil {
				log.Warn("daemon.Daemon.run error: %v", err)
				break
			}
		}
		d.ready()
		<-ctx.Done()
	}
	return err
}

func (d *Daemon) init() error {
	return d.Start(d.cfg.Username, d.cfg.Password)
}

func (d *Daemon) Start(username, password string) error {
	resp, err := d.athm.SignIn(username, password)
	if err != nil {
		return fmt.Errorf("auth sign in, error: %v", err)
	}
	d.uin = resp.Data.Uin
	if err := d.sync(d.uin); err != nil {
		return err
	}
	if _, err := d.svcm.RegisterSetOnlineStatus(d.uin, service.StatusTypeOnline, true); err != nil {
		return fmt.Errorf("service register set online status, error: %v", err)
	}
	d.started()
	log.Info("daemon is started, start handling push message")
	go d.watchDog(d.uin)
	return nil
}

func (d *Daemon) sync(uin int64) error {
	if _, err := d.svcm.RegisterAppRegister(d.uin); err != nil {
		return fmt.Errorf("service register app register, error: %v", err)
	}
	if _, err := d.cntm.GetContacts(d.uin, 0, 100, 0, 100); err != nil {
		return fmt.Errorf("contact get contacts, error: %v", err)
	}
	if _, err := d.chtm.GetGroups(d.uin); err != nil {
		return fmt.Errorf("chat get groups and users, error: %v", err)
	}
	if _, err := d.chnm.SyncFirstView(d.uin, 0); err != nil {
		return fmt.Errorf("channel sync first view, error: %v", err)
	}
	return nil
}

func (d *Daemon) resync(uin int64) error {
	if _, err := d.svcm.RegisterAppRegister(d.uin); err != nil {
		return fmt.Errorf("service register app register, error: %v", err)
	}
	return nil
}

func (d *Daemon) ready() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.isReady {
		d.isReady = true
		if len(d.readyCh) == 0 {
			d.readyCh <- struct{}{}
		}
	}
}

func (d *Daemon) checkReady() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isReady
}

func (d *Daemon) Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error {
	err := d.c.Call(serviceMethod, args, reply)
	if err != nil {
		d.mu.Lock()
		d.lastRecvTime = time.Now()
		d.mu.Unlock()
	}
	return err
}

func (d *Daemon) started() {
	d.mu.Lock()
	if d.waiting {
		d.wg.Done()
		d.waiting = false
	} else {
		log.Warn("daemon already started")
	}
	d.mu.Unlock()
}

func (d *Daemon) Waiting() {
	d.mu.Lock()
	if !d.waiting {
		d.wg.Add(1)
		d.waiting = true
	} else {
		log.Warn("daemon already waiting")
	}
	d.mu.Unlock()
}

func (d *Daemon) Register(serviceMethod string, handler rpc.Handler) error {
	return d.c.Register(serviceMethod, func(reply *rpc.Reply) (*rpc.Args, error) {
		{
			d.mu.Lock()
			if d.waiting {
				log.Warn("daemon is waiting, pending push %s", serviceMethod)
			}
			d.mu.Unlock()
		}
		d.wg.Wait()
		args, err := handler(reply)
		if err != nil {
			return nil, err
		}
		if args != nil {
			d.mu.Lock()
			d.lastRecvTime = time.Now()
			d.mu.Unlock()
		}
		return args, nil
	})
}

func (d *Daemon) GetFakeSource(uin int64) *fake.Source {
	return d.c.GetFakeSource(uin)
}

func (d *Daemon) GetMessageChannel() chan *penguin.Message {
	return d.msgChan
}

func (d *Daemon) GetAccountManager() *account.Manager {
	return d.accm
}

func (d *Daemon) GetAuthManager() *auth.Manager {
	return d.athm
}

func (d *Daemon) GetChannelManager() *channel.Manager {
	return d.chnm
}

func (d *Daemon) GetChatManager() *chat.Manager {
	return d.chtm
}

func (d *Daemon) GetContactManager() *contact.Manager {
	return d.cntm
}

func (d *Daemon) GetMessageManager() *message.Manager {
	return d.msgm
}

func (d *Daemon) GetServiceManager() *service.Manager {
	return d.svcm
}
