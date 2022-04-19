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

package message

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Manager struct {
	Daemon
	ctx context.Context

	mu sync.RWMutex
	// session
	flags   map[int64]int32
	cookies map[int64][]byte
}

func NewManager(ctx context.Context, d Daemon) *Manager {
	m := &Manager{
		Daemon: d,
		ctx:    ctx,
		flags:  make(map[int64]int32),
	}
	m.Register(service.MethodMessagePushNotify, m.handlePushNotifyRequest)
	return m
}

func (m *Manager) GetFlag(k int64) (int32, bool) {
	m.mu.RLock()
	v, ok := m.flags[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetFlag(k int64, v int32) (int32, bool) {
	m.mu.Lock()
	vv, ok := m.flags[k]
	m.flags[k] = v
	m.mu.Unlock()
	return vv, ok
}

func (m *Manager) GetCookie(k int64) ([]byte, bool) {
	m.mu.RLock()
	if m.cookies == nil {
		m.cookies = getCookies()
	}
	v, ok := m.cookies[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *Manager) SetCookie(k int64, v []byte) ([]byte, bool) {
	m.mu.Lock()
	vv, ok := m.cookies[k]
	m.cookies[k] = v
	setCookies(m.cookies)
	m.mu.Unlock()
	return vv, ok
}

func getCookies() map[int64][]byte {
	home, _ := os.UserHomeDir()
	file := path.Join(home, ".penguin/session/message.cookies.json")
	data, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		data, err = json.MarshalIndent(map[int64][]byte{}, "", "  ")
		if err == nil {
			err = ioutil.WriteFile(file, data, 0644)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}
	var cookies map[int64][]byte
	err = json.Unmarshal(data, &cookies)
	if err != nil {
		log.Fatalln(err)
	}
	return cookies
}

func setCookies(cookies map[int64][]byte) {
	home, _ := os.UserHomeDir()
	file := path.Join(home, ".penguin/session/message.cookies.json")
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(file, data, 0644)
	}
}
