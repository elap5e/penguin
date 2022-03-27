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

	"github.com/elap5e/penguin/config"
	"github.com/elap5e/penguin/daemon/account"
	"github.com/elap5e/penguin/daemon/auth"
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
	cntm *contact.Manager
	msgm *message.Manager
	svcm *service.Manager
}

func New(ctx context.Context, cfg *config.Config) *Daemon {
	d := &Daemon{
		ctx: ctx,
		cfg: cfg,
		c:   msf.NewClient(ctx),
	}
	d.athm = auth.NewManager(d.ctx, d.c)
	d.svcm = service.NewManager(d.ctx, d.c, d)
	return d
}

func (d *Daemon) Run() error {
	resp, err := d.athm.SignIn(d.cfg.Username, d.cfg.Password)
	if err != nil {
		return fmt.Errorf("sign in, error: %v", err)
	}
	if _, err := d.svcm.RegisterAppRegister(resp.Data.Uin); err != nil {
		return fmt.Errorf("register app register, error: %v", err)
	}
	select {}
}

func (d *Daemon) GetAuthManager() *auth.Manager {
	return d.athm
}
