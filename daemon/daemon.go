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
	"encoding/json"
	"log"

	"github.com/elap5e/penguin/daemon/account"
	"github.com/elap5e/penguin/daemon/auth"
	"github.com/elap5e/penguin/daemon/contact"
	"github.com/elap5e/penguin/daemon/message"
	"github.com/elap5e/penguin/pkg/net/msf"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Daemon struct {
	ctx context.Context

	c rpc.Client

	accm *account.Manager
	athm *auth.Manager
	cntm *contact.Manager
	msgm *message.Manager
}

func New(ctx context.Context) *Daemon {
	d := &Daemon{
		ctx: ctx,
		c:   msf.NewClient(ctx),
	}
	d.athm = auth.NewManager(d.ctx, d.c)
	return d
}

func (d *Daemon) Run() error {
	var p []byte
	call := <-d.c.Go(service.MethodHeartbeatAlive, &rpc.Args{
		Uin:     0,
		Seq:     d.c.GetNextSeq(),
		Payload: []byte{0, 0, 0, 4},
	}, &rpc.Reply{}, make(chan *rpc.Call, 1)).Done
	p, _ = json.MarshalIndent(call.Reply, "", "  ")
	log.Printf("call.Reply:\n%s", string(p))
	_, err := d.athm.SignIn("10000", "password")
	if err != nil {
		return err
	}
	return nil
}
