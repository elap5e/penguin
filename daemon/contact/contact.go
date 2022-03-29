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

package contact

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

	mu       sync.RWMutex
	contacts map[int64]map[int64]*penguin.Contact
}

func NewManager(ctx context.Context, c rpc.Client, d Daemon) *Manager {
	return &Manager{
		ctx:      ctx,
		c:        c,
		d:        d,
		contacts: make(map[int64]map[int64]*penguin.Contact),
	}
}

func (m *Manager) getContacts(uin int64) map[int64]*penguin.Contact {
	contacts, ok := m.contacts[uin]
	if !ok {
		m.contacts[uin] = make(map[int64]*penguin.Contact)
		contacts = m.contacts[uin]
	}
	return contacts
}

func (m *Manager) GetContact(uin, k int64) (*penguin.Contact, bool) {
	m.mu.RLock()
	contacts := m.getContacts(uin)
	v, ok := contacts[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetContact(uin, k int64, v *penguin.Contact) (*penguin.Contact, bool) {
	m.mu.Lock()
	contacts := m.getContacts(uin)
	vv, ok := contacts[k]
	contacts[k] = v
	m.mu.Unlock()
	return vv, ok
}
