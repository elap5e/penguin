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
	"fmt"
	"math/rand"
	"sync"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/pkg/log"
)

type Store interface {
	GetChat(chatID int64) (oldChat *penguin.Chat, ok bool)
	SetChat(chatID int64, newChat *penguin.Chat) (oldChat *penguin.Chat, ok bool)

	GetChatUser(chatID, userID int64) (oldUser *penguin.User, ok bool)
	SetChatUser(chatID, userID int64, newUser *penguin.User) (oldUser *penguin.User, ok bool)

	GetChatSeq(accountID, chatID, userID int64) (oldSeq uint32, ok bool)
	GetNextChatSeq(accountID, chatID, userID int64) (seq uint32, ok bool)
	SetChatSeq(accountID, chatID, userID int64, newSeq uint32) (oldSeq uint32, ok bool)
	GetCookie(userID int64) (oldCookie []byte, ok bool)
	SetCookie(userID int64, newCookie []byte) (oldCookie []byte, ok bool)
}

type memStore struct {
	mu    sync.RWMutex
	chats map[int64]*penguin.Chat           // shared
	users map[int64]map[int64]*penguin.User // shared

	// session
	chatSeq map[string]uint32
	cookies map[int64][]byte
}

func NewMemStore() Store {
	return &memStore{
		mu:      sync.RWMutex{},
		chats:   make(map[int64]*penguin.Chat),
		users:   make(map[int64]map[int64]*penguin.User),
		chatSeq: make(map[string]uint32),
		cookies: make(map[int64][]byte),
	}
}

func copyChat(newChat, oldChat *penguin.Chat) {
	if oldChat == nil {
		log.Warn("chat.copyChat: oldChat is nil")
		return
	} else if newChat == nil {
		log.Warn("chat.copyChat: newChat is nil")
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

func (s *memStore) GetChat(chatID int64) (oldChat *penguin.Chat, ok bool) {
	s.mu.RLock()
	chat, ok := s.chats[chatID]
	if ok {
		oldChat = new(penguin.Chat)
		copyChat(oldChat, chat)
	}
	s.mu.RUnlock()
	return oldChat, ok
}

func (s *memStore) SetChat(chatID int64, newChat *penguin.Chat) (oldChat *penguin.Chat, ok bool) {
	s.mu.Lock()
	chat, ok := s.chats[chatID]
	if !ok {
		s.chats[chatID] = new(penguin.Chat)
		chat = s.chats[chatID]
	} else {
		oldChat = new(penguin.Chat)
		copyChat(oldChat, chat)
	}
	copyChat(chat, newChat)
	s.mu.Unlock()
	return oldChat, ok
}

func copyUser(newChatUser, oldChatUser *penguin.User) {
	if oldChatUser == nil {
		log.Warn("chat.copyUser: oldChatUser is nil")
		return
	} else if newChatUser == nil {
		log.Warn("chat.copyUser: newChatUser is nil")
		return
	}
	newChatUser.Account = oldChatUser.Account
	newChatUser.Display = oldChatUser.Display
	return
}

func (s *memStore) getUsersByChatID(chatID int64) map[int64]*penguin.User {
	s.mu.Lock()
	users, ok := s.users[chatID]
	if !ok {
		s.users[chatID] = make(map[int64]*penguin.User)
		users = s.users[chatID]
	}
	s.mu.Unlock()
	return users
}

func (s *memStore) GetChatUser(chatID, userID int64) (oldUser *penguin.User, ok bool) {
	users := s.getUsersByChatID(chatID)
	s.mu.RLock()
	user, ok := users[userID]
	if ok {
		oldUser = new(penguin.User)
		copyUser(oldUser, user)
	}
	s.mu.RUnlock()
	return oldUser, ok
}

func (s *memStore) SetChatUser(chatID, userID int64, newUser *penguin.User) (oldUser *penguin.User, ok bool) {
	users := s.getUsersByChatID(chatID)
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

func (s *memStore) GetChatSeq(accountID, chatID, userID int64) (oldSeq uint32, ok bool) {
	s.mu.RLock()
	key := fmt.Sprintf("%d|%d|%d", accountID, chatID, userID)
	oldSeq, ok = s.chatSeq[key]
	s.mu.RUnlock()
	return oldSeq, ok
}

func (s *memStore) GetNextChatSeq(accountID, chatID, userID int64) (seq uint32, ok bool) {
	s.mu.RLock()
	key := fmt.Sprintf("%d|%d|%d", accountID, chatID, userID)
	oldSeq, ok := s.chatSeq[key]
	if oldSeq == 0 {
		oldSeq = uint32(rand.Int31n(100000) + 10000)
	} else if oldSeq > 1000000 {
		oldSeq = uint32(rand.Int31n(100000) + 60000)
	}
	s.chatSeq[key] = oldSeq + 1
	s.mu.RUnlock()
	return oldSeq + 1, ok
}

func (s *memStore) SetChatSeq(accountID, chatID, userID int64, newSeq uint32) (oldSeq uint32, ok bool) {
	s.mu.Lock()
	key := fmt.Sprintf("%d|%d|%d", accountID, chatID, userID)
	oldSeq, ok = s.chatSeq[key]
	s.chatSeq[key] = newSeq
	s.mu.Unlock()
	return oldSeq, ok
}

func (s *memStore) GetCookie(userID int64) (oldCookie []byte, ok bool) {
	s.mu.RLock()
	cookie, ok := s.cookies[userID]
	oldCookie = append(oldCookie, cookie...)
	s.mu.RUnlock()
	return oldCookie, ok
}

func (s *memStore) SetCookie(userID int64, newCookie []byte) (oldCookie []byte, ok bool) {
	s.mu.Lock()
	cookie, ok := s.cookies[userID]
	oldCookie = append(oldCookie, cookie...)
	s.cookies[userID] = newCookie
	s.mu.Unlock()
	return oldCookie, ok
}
