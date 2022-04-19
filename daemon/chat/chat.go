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

package chat

import (
	"context"

	"github.com/elap5e/penguin/daemon/account"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Daemon interface {
	Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error

	GetAccountManager() *account.Manager
}

type Manager struct {
	context.Context
	Daemon
	Store
}

func NewManager(ctx context.Context, daemon Daemon, store Store) *Manager {
	return &Manager{
		Context: ctx,
		Daemon:  daemon,
		Store:   store,
	}
}
