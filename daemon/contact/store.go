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
	"sync"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/pkg/log"
)

type Store interface {
	GetContact(accountID, contactID int64) (oldContact *penguin.Contact, ok bool)
	SetContact(accountID, contactID int64, newContact *penguin.Contact) (oldContact *penguin.Contact, ok bool)
}

type memStore struct {
	mu       sync.RWMutex
	contacts map[int64]map[int64]*penguin.Contact
}

func NewMemStore() Store {
	return &memStore{
		mu:       sync.RWMutex{},
		contacts: make(map[int64]map[int64]*penguin.Contact),
	}
}

func (s *memStore) getContactsByAccountID(accountID int64) map[int64]*penguin.Contact {
	s.mu.Lock()
	contacts, ok := s.contacts[accountID]
	if !ok {
		s.contacts[accountID] = make(map[int64]*penguin.Contact)
		contacts = s.contacts[accountID]
	}
	s.mu.Unlock()
	return contacts
}

func copyContact(newContact, oldContact *penguin.Contact) {
	if oldContact == nil {
		log.Warn("contact.copyContact: oldContact is nil")
		return
	} else if newContact == nil {
		log.Warn("contact.copyContact: newContact is nil")
		return
	}
	newContact.Account = oldContact.Account
	newContact.Display = oldContact.Display
	return
}

func (s *memStore) GetContact(accountID, contactID int64) (oldContact *penguin.Contact, ok bool) {
	contacts := s.getContactsByAccountID(accountID)
	s.mu.RLock()
	contact, ok := contacts[contactID]
	if ok {
		oldContact = new(penguin.Contact)
		copyContact(oldContact, contact)
	}
	s.mu.RUnlock()
	return oldContact, ok
}

func (s *memStore) SetContact(accountID, contactID int64, newContact *penguin.Contact) (oldContact *penguin.Contact, ok bool) {
	contacts := s.getContactsByAccountID(accountID)
	s.mu.Lock()
	contact, ok := contacts[contactID]
	if !ok {
		contacts[contactID] = new(penguin.Contact)
		contact = contacts[contactID]
	} else {
		oldContact = new(penguin.Contact)
		copyContact(oldContact, contact)
	}
	copyContact(contact, newContact)
	s.mu.Unlock()
	return contact, ok
}
