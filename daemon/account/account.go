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
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Daemon interface {
	Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error
}

type Manager struct {
	Daemon
	ctx context.Context

	mu       sync.RWMutex
	defaults map[int64]*penguin.Account // shared
	channels map[int64]*penguin.Account // shared
}

func NewManager(ctx context.Context, d Daemon) *Manager {
	return &Manager{
		Daemon:   d,
		ctx:      ctx,
		defaults: make(map[int64]*penguin.Account),
		channels: make(map[int64]*penguin.Account),
	}
}

func (m *Manager) GetDefaultAccount(k int64) (*penguin.Account, bool) {
	m.mu.RLock()
	v, ok := m.defaults[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetDefaultAccount(k int64, v *penguin.Account) (*penguin.Account, bool) {
	m.mu.Lock()
	vv, ok := m.defaults[k]
	m.defaults[k] = v
	m.mu.Unlock()
	return vv, ok
}

func (m *Manager) GetChannelAccount(k int64) (*penguin.Account, bool) {
	m.mu.RLock()
	v, ok := m.channels[k]
	m.mu.RUnlock()
	if v.ID>>32 == 0 {
		v.ID |= 0x0200000700000000
	}
	return v, ok
}

func (m *Manager) SetChannelAccount(k int64, v *penguin.Account) (*penguin.Account, bool) {
	m.mu.Lock()
	vv, ok := m.channels[k]
	if v.ID>>32 != 0x02000007 {
		log.Warn("account.Manager.SetChannelAccount: invalid account id: 0x%x", v.ID)
	} else {
		v.ID &= 0x00000000ffffffff
	}
	m.channels[k] = v
	m.mu.Unlock()
	return vv, ok
}
