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
	"sync"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/pkg/log"
)

type Store interface {
	GetDefaultAccount(id int64) (account *penguin.Account, ok bool)
	SetDefaultAccount(id int64, newAccount *penguin.Account) (oldAccount *penguin.Account, ok bool)

	GetChannelAccount(id int64) (account *penguin.Account, ok bool)
	SetChannelAccount(id int64, newAccount *penguin.Account) (oldAccount *penguin.Account, ok bool)

	GetChannelFromDefault(id int64) (*penguin.Account, bool)
	GetDefaultFromChannel(id int64) (*penguin.Account, bool)
	SetDefaultChannelPair(src, dst int64)
}

type memStore struct {
	mu       sync.RWMutex
	defaults map[int64]*penguin.Account // shared
	channels map[int64]*penguin.Account // shared

	defaultToChannel map[int64]int64 // shared
	channelToDefault map[int64]int64 // shared
}

func NewMemStore() Store {
	return &memStore{
		mu:       sync.RWMutex{},
		defaults: make(map[int64]*penguin.Account),
		channels: make(map[int64]*penguin.Account),

		defaultToChannel: make(map[int64]int64),
		channelToDefault: make(map[int64]int64),
	}
}

func copyDefaultAccount(newAccount, oldAccount *penguin.Account) {
	if oldAccount == nil {
		log.Warn("account.copyDefaultAccount: oldAccount is nil")
		return
	} else if newAccount == nil {
		log.Warn("account.copyDefaultAccount: newAccount is nil")
		return
	}
	newAccount.ID = oldAccount.ID
	newAccount.Type = oldAccount.Type
	newAccount.Username = oldAccount.Username
	newAccount.Photo = oldAccount.Photo
	return
}

func (s *memStore) GetDefaultAccount(id int64) (oldAccount *penguin.Account, ok bool) {
	s.mu.RLock()
	account, ok := s.defaults[id]
	if ok {
		oldAccount = new(penguin.Account)
		copyDefaultAccount(oldAccount, account)
	}
	s.mu.RUnlock()
	return oldAccount, ok
}

func (s *memStore) SetDefaultAccount(id int64, newAccount *penguin.Account) (oldAccount *penguin.Account, ok bool) {
	if id>>32 != 0 {
		log.Warn("account.SetDefaultAccount: invalid account id: 0x%x", id)
	}
	s.mu.Lock()
	account, ok := s.defaults[id]
	if !ok {
		s.defaults[id] = new(penguin.Account)
		account = s.defaults[id]
	} else {
		oldAccount = new(penguin.Account)
		copyDefaultAccount(oldAccount, account)
	}
	copyDefaultAccount(account, newAccount)
	s.mu.Unlock()
	return oldAccount, ok
}

func copyChannelAccount(newAccount, oldAccount *penguin.Account, set ...bool) {
	if oldAccount == nil {
		log.Warn("account.copyChannelAccount: oldAccount is nil")
		return
	} else if newAccount == nil {
		log.Warn("account.copyChannelAccount: newAccount is nil")
		return
	}
	if len(set) > 0 && set[0] {
		newAccount.ID = setChannelAccountID(oldAccount.ID)
	} else {
		newAccount.ID = getChannelAccountID(oldAccount.ID)
	}
	newAccount.Type = oldAccount.Type
	newAccount.Username = oldAccount.Username
	newAccount.Photo = oldAccount.Photo
	return
}

func getChannelAccountID(id int64) int64 {
	if id>>32 == 0 {
		return id | 0x0200000700000000
	}
	log.Warn("account.getChannelAccountID: invalid account id: 0x%x", id)
	return id
}

func setChannelAccountID(id int64) int64 {
	if id>>32 == 0x02000007 {
		return id & 0x00000000ffffffff
	}
	log.Warn("account.setChannelAccountID: invalid account id: 0x%x", id)
	return id
}

func (s *memStore) GetChannelAccount(id int64) (oldAccount *penguin.Account, ok bool) {
	s.mu.RLock()
	account, ok := s.defaults[setChannelAccountID(id)]
	if ok {
		oldAccount = new(penguin.Account)
		copyChannelAccount(oldAccount, account)
	}
	s.mu.RUnlock()
	return oldAccount, ok
}

func (s *memStore) SetChannelAccount(id int64, newAccount *penguin.Account) (oldAccount *penguin.Account, ok bool) {
	s.mu.Lock()
	id = setChannelAccountID(id)
	account, ok := s.defaults[id]
	if !ok {
		s.defaults[id] = new(penguin.Account)
		account = s.defaults[id]
	} else {
		oldAccount = new(penguin.Account)
		copyChannelAccount(oldAccount, account)
	}
	copyChannelAccount(account, newAccount, true)
	s.mu.Unlock()
	return oldAccount, ok
}

func (s *memStore) GetChannelFromDefault(id int64) (*penguin.Account, bool) {
	s.mu.RLock()
	id, ok := s.defaultToChannel[id]
	s.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return s.GetChannelAccount(id)
}

func (s *memStore) GetDefaultFromChannel(id int64) (*penguin.Account, bool) {
	s.mu.RLock()
	id, ok := s.channelToDefault[id]
	s.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return s.GetDefaultAccount(id)
}

func (s *memStore) SetDefaultChannelPair(src, dst int64) {
	s.mu.RLock()
	s.channelToDefault[dst] = src
	s.defaultToChannel[src] = dst
	s.mu.RUnlock()
}
