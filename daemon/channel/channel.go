package channel

import (
	"context"
	"net/rpc"
	"sync"

	"github.com/elap5e/penguin"
)

type Daemon interface {
}

type Manager struct {
	ctx context.Context

	c rpc.Client
	d Daemon

	mu       sync.RWMutex
	channels map[int64]*penguin.Chat // shared
	rooms    map[int64]*penguin.Chat // shared
	users    map[int64]map[int64]*penguin.User
}

func NewManager(ctx context.Context, c rpc.Client, d Daemon) *Manager {
	return &Manager{
		ctx:      ctx,
		c:        c,
		d:        d,
		channels: make(map[int64]*penguin.Chat),
		rooms:    make(map[int64]*penguin.Chat),
		users:    make(map[int64]map[int64]*penguin.User),
	}
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

func (m *Manager) GetRoom(k int64) (*penguin.Chat, bool) {
	m.mu.RLock()
	v, ok := m.rooms[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetRoom(k int64, v *penguin.Chat) (*penguin.Chat, bool) {
	m.mu.Lock()
	vv, ok := m.rooms[k]
	m.rooms[k] = v
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
