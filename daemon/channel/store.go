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
	"sync"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/pkg/log"
)

type Store interface {
	GetChannel(channelID int64) (oldChannel *penguin.Chat, ok bool)
	SetChannel(channelID int64, newChannel *penguin.Chat) (oldChannel *penguin.Chat, ok bool)

	GetRoom(channelID, roomID int64) (oldRoom *penguin.Chat, ok bool)
	SetRoom(channelID, roomID int64, newRoom *penguin.Chat) (oldRoom *penguin.Chat, ok bool)

	GetUser(channelID, userID int64) (oldUser *penguin.User, ok bool)
	SetUser(channelID, userID int64, newUser *penguin.User) (oldUser *penguin.User, ok bool)
}

type memStore struct {
	mu       sync.RWMutex
	channels map[int64]*penguin.Chat           // shared
	roles    map[int64]map[int64]*penguin.Role // shared
	rooms    map[int64]map[int64]*penguin.Chat // shared
	users    map[int64]map[int64]*penguin.User
}

func NewMemStore() Store {
	return &memStore{
		mu:       sync.RWMutex{},
		channels: make(map[int64]*penguin.Chat),
		roles:    make(map[int64]map[int64]*penguin.Role),
		rooms:    make(map[int64]map[int64]*penguin.Chat),
		users:    make(map[int64]map[int64]*penguin.User),
	}
}

func copyChat(newChat, oldChat *penguin.Chat) {
	if oldChat == nil {
		log.Warn("channel.copyChat: oldChat is nil")
		return
	} else if newChat == nil {
		log.Warn("channel.copyChat: newChat is nil")
		return
	}
	newChat.ID = oldChat.ID
	newChat.Type = oldChat.Type
	newChat.Chat = oldChat.Chat
	newChat.Channel = oldChat.Channel
	newChat.User = oldChat.User
	newChat.Title = oldChat.Title
	newChat.Photo = oldChat.Photo
	newChat.Display = oldChat.Display
	return
}

func (s *memStore) GetChannel(channelID int64) (oldChannel *penguin.Chat, ok bool) {
	s.mu.RLock()
	channel, ok := s.channels[channelID]
	if ok {
		oldChannel = new(penguin.Chat)
		copyChat(oldChannel, channel)
	}
	s.mu.RUnlock()
	return oldChannel, ok
}

func (s *memStore) SetChannel(channelID int64, newChannel *penguin.Chat) (oldChannel *penguin.Chat, ok bool) {
	s.mu.Lock()
	channel, ok := s.channels[channelID]
	if !ok {
		s.channels[channelID] = new(penguin.Chat)
		channel = s.channels[channelID]
	} else {
		oldChannel = new(penguin.Chat)
		copyChat(oldChannel, channel)
	}
	copyChat(channel, newChannel)
	s.mu.Unlock()
	return oldChannel, ok
}

func (s *memStore) getRoomsByChannelID(channelID int64) map[int64]*penguin.Chat {
	s.mu.Lock()
	rooms, ok := s.rooms[channelID]
	if !ok {
		s.rooms[channelID] = make(map[int64]*penguin.Chat)
		rooms = s.rooms[channelID]
	}
	s.mu.Unlock()
	return rooms
}

func (s *memStore) GetRoom(channelID, roomID int64) (oldRoom *penguin.Chat, ok bool) {
	rooms := s.getRoomsByChannelID(channelID)
	s.mu.RLock()
	room, ok := rooms[roomID]
	if ok {
		oldRoom = new(penguin.Chat)
		copyChat(oldRoom, room)
	}
	s.mu.RUnlock()
	return oldRoom, ok
}

func (s *memStore) SetRoom(channelID, roomID int64, newRoom *penguin.Chat) (oldRoom *penguin.Chat, ok bool) {
	rooms := s.getRoomsByChannelID(channelID)
	s.mu.Lock()
	room, ok := rooms[roomID]
	if !ok {
		rooms[roomID] = new(penguin.Chat)
		room = rooms[roomID]
	} else {
		oldRoom = new(penguin.Chat)
		copyChat(oldRoom, room)
	}
	copyChat(room, newRoom)
	s.mu.Unlock()
	return oldRoom, ok
}

func (s *memStore) getRolesByChannelID(channelID int64) map[int64]*penguin.Role {
	s.mu.Lock()
	roles, ok := s.roles[channelID]
	if !ok {
		s.roles[channelID] = make(map[int64]*penguin.Role)
		roles = s.roles[channelID]
	}
	s.mu.Unlock()
	return roles
}

func (s *memStore) GetRole(channelID, roleID int64) (oldRole *penguin.Role, ok bool) {
	roles := s.getRolesByChannelID(roleID)
	s.mu.RLock()
	v, ok := roles[roleID]
	s.mu.RUnlock()
	return v, ok
}

func (s *memStore) SetRole(channelID, roleID int64, newRole *penguin.Role) (oldRole *penguin.Role, ok bool) {
	roles := s.getRolesByChannelID(channelID)
	s.mu.Lock()
	vv, ok := roles[roleID]
	roles[roleID] = newRole
	s.mu.Unlock()
	return vv, ok
}

func copyUser(newUser, oldUser *penguin.User) {
	if oldUser == nil {
		log.Warn("channel.copyUser: oldUser is nil")
		return
	} else if newUser == nil {
		log.Warn("channel.copyUser: newUser is nil")
		return
	}
	newUser.Account = oldUser.Account
	newUser.Display = oldUser.Display
	return
}

func (s *memStore) getUsersByChannelID(channelID int64) map[int64]*penguin.User {
	s.mu.Lock()
	users, ok := s.users[channelID]
	if !ok {
		s.users[channelID] = make(map[int64]*penguin.User)
		users = s.users[channelID]
	}
	s.mu.Unlock()
	return users
}

func (s *memStore) GetUser(channelID, userID int64) (oldUser *penguin.User, ok bool) {
	users := s.getUsersByChannelID(channelID)
	s.mu.RLock()
	user, ok := users[userID]
	if ok {
		oldUser = new(penguin.User)
		copyUser(oldUser, user)
	}
	s.mu.RUnlock()
	return oldUser, ok
}

func (s *memStore) SetUser(channelID, userID int64, newUser *penguin.User) (oldUser *penguin.User, ok bool) {
	users := s.getUsersByChannelID(channelID)
	s.mu.Lock()
	user, ok := users[userID]
	if !ok {
		users[userID] = new(penguin.User)
		user = users[userID]
	} else {
		oldUser = new(penguin.User)
		copyUser(oldUser, user)
	}
	copyUser(user, newUser)
	s.mu.Unlock()
	return oldUser, ok
}
