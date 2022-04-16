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
	"github.com/elap5e/penguin/pkg/net/msf"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Daemon struct {
	ctx context.Context
	cfg *config.Config

	c rpc.Client

	accm *account.Manager
	athm *auth.Manager
	chtm *chat.Manager
	chnm *channel.Manager
	cntm *contact.Manager
	msgm *message.Manager
	svcm *service.Manager

	msgChan chan *penguin.Message
}

func New(ctx context.Context, cfg *config.Config) *Daemon {
	d := &Daemon{
		ctx: ctx,
		cfg: cfg,
		c:   msf.NewClient(ctx),
	}
	d.accm = account.NewManager(d.ctx, d.c)
	d.athm = auth.NewManager(d.ctx, d.c)
	d.chtm = chat.NewManager(d.ctx, d.c, d)
	d.chnm = channel.NewManager(d.ctx, d.c, d)
	d.cntm = contact.NewManager(d.ctx, d.c, d)
	d.msgm = message.NewManager(d.ctx, d.c, d)
	d.svcm = service.NewManager(d.ctx, d.c, d)
	d.msgChan = make(chan *penguin.Message, 100)
	return d
}

func (d *Daemon) Run() error {
	resp, err := d.athm.SignIn(d.cfg.Username, d.cfg.Password)
	if err != nil {
		return fmt.Errorf("auth sign in, error: %v", err)
	}
	if _, err := d.svcm.RegisterAppRegister(resp.Data.Uin); err != nil {
		return fmt.Errorf("service register app register, error: %v", err)
	}
	time.Sleep(time.Second * 5)
	if _, err := d.cntm.GetContacts(resp.Data.Uin, 0, 100, 0, 100); err != nil {
		return fmt.Errorf("contact get contacts, error: %v", err)
	}
	if _, err := d.chtm.GetGroups(resp.Data.Uin); err != nil {
		return fmt.Errorf("chat get groups and users, error: %v", err)
	}
	if _, err := d.chnm.SyncFirstView(resp.Data.Uin, 0); err != nil {
		return fmt.Errorf("channel sync first view, error: %v", err)
	}
	if _, err := d.svcm.RegisterSetOnlineStatus(resp.Data.Uin, service.StatusTypeOnline, true); err != nil {
		return fmt.Errorf("service register set online status, error: %v", err)
	}
	select {}
}

func (d *Daemon) GetAccountManager() *account.Manager {
	return d.accm
}

func (d *Daemon) GetAuthManager() *auth.Manager {
	return d.athm
}

func (d *Daemon) GetServiceManager() *service.Manager {
	return d.svcm
}

func (d *Daemon) GetMessageChannel() chan *penguin.Message {
	return d.msgChan
}
