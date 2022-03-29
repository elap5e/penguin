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
	"sync"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/account"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Daemon interface {
	GetAccountManager() *account.Manager
}

type Manager struct {
	ctx context.Context

	c rpc.Client
	d Daemon

	mu    sync.RWMutex
	chats map[int64]*penguin.Chat           // shared
	users map[int64]map[int64]*penguin.User // shared

	// session
	cookies map[int64][]byte
}

func NewManager(ctx context.Context, c rpc.Client, d Daemon) *Manager {
	return &Manager{
		ctx:     ctx,
		c:       c,
		d:       d,
		chats:   make(map[int64]*penguin.Chat),
		users:   make(map[int64]map[int64]*penguin.User),
		cookies: make(map[int64][]byte),
	}
}

func (m *Manager) GetChat(k int64) (*penguin.Chat, bool) {
	m.mu.RLock()
	v, ok := m.chats[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetChat(k int64, v *penguin.Chat) (*penguin.Chat, bool) {
	m.mu.Lock()
	vv, ok := m.chats[k]
	m.chats[k] = v
	m.mu.Unlock()
	return vv, ok
}

func (m *Manager) getUsers(uin int64) map[int64]*penguin.User {
	users, ok := m.users[uin]
	if !ok {
		m.users[uin] = make(map[int64]*penguin.User)
		users = m.users[uin]
	}
	return users
}

func (m *Manager) GetUser(uin, k int64) (*penguin.User, bool) {
	m.mu.RLock()
	users := m.getUsers(uin)
	v, ok := users[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetUser(uin, k int64, v *penguin.User) (*penguin.User, bool) {
	m.mu.Lock()
	users := m.getUsers(uin)
	vv, ok := users[k]
	users[k] = v
	m.mu.Unlock()
	return vv, ok
}

func (m *Manager) GetCookie(k int64) ([]byte, bool) {
	m.mu.RLock()
	v, ok := m.cookies[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetCookie(k int64, v []byte) ([]byte, bool) {
	m.mu.Lock()
	vv, ok := m.cookies[k]
	m.cookies[k] = v
	m.mu.Unlock()
	return vv, ok
}
