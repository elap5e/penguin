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

package account

import (
	"context"
	"sync"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Daemon interface {
	Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error
}

type Manager struct {
	Daemon
	ctx context.Context

	mu       sync.RWMutex
	accounts map[int64]*penguin.Account // shared
}

func NewManager(ctx context.Context, d Daemon) *Manager {
	return &Manager{
		Daemon:   d,
		ctx:      ctx,
		accounts: make(map[int64]*penguin.Account),
	}
}

func (m *Manager) GetAccount(k int64) (*penguin.Account, bool) {
	m.mu.RLock()
	v, ok := m.accounts[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetAccount(k int64, v *penguin.Account) (*penguin.Account, bool) {
	m.mu.Lock()
	vv, ok := m.accounts[k]
	m.accounts[k] = v
	m.mu.Unlock()
	return vv, ok
}
