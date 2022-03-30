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

package channel

import (
	"context"
	"sync"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/channel/pb"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Daemon interface {
	OnRecvChannelMessage(id int64, recv *pb.Common_Msg) error
}

type Manager struct {
	ctx context.Context

	c rpc.Client
	d Daemon

	mu       sync.RWMutex
	channels map[int64]*penguin.Chat           // shared
	rooms    map[int64]map[int64]*penguin.Chat // shared
	users    map[int64]map[int64]*penguin.User
}

func NewManager(ctx context.Context, c rpc.Client, d Daemon) *Manager {
	m := &Manager{
		ctx:      ctx,
		c:        c,
		d:        d,
		channels: make(map[int64]*penguin.Chat),
		rooms:    make(map[int64]map[int64]*penguin.Chat),
		users:    make(map[int64]map[int64]*penguin.User),
	}
	m.c.Register(service.MethodChannelPushMessage, m.handlePushMessage)
	m.c.Register(service.MethodChannelPushFirstView, m.handlePushFirstView)
	return m
}

func (m *Manager) GetChannel(k int64) (*penguin.Chat, bool) {
	m.mu.RLock()
	v, ok := m.channels[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetChannel(k int64, v *penguin.Chat) (*penguin.Chat, bool) {
	m.mu.Lock()
	vv, ok := m.channels[k]
	m.channels[k] = v
	m.mu.Unlock()
	return vv, ok
}

func (m *Manager) getRooms(cid int64) map[int64]*penguin.Chat {
	rooms, ok := m.rooms[cid]
	if !ok {
		m.rooms[cid] = make(map[int64]*penguin.Chat)
		rooms = m.rooms[cid]
	}
	return rooms
}

func (m *Manager) GetRoom(cid, rid int64) (*penguin.Chat, bool) {
	m.mu.RLock()
	rooms := m.getRooms(rid)
	v, ok := rooms[rid]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetRoom(cid, rid int64, v *penguin.Chat) (*penguin.Chat, bool) {
	m.mu.Lock()
	rooms := m.getRooms(cid)
	vv, ok := rooms[rid]
	rooms[rid] = v
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
