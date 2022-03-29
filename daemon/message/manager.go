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

	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Manager struct {
	ctx context.Context

	c rpc.Client
	d Daemon

	syncFlags   map[int64]int32
	syncCookies map[int64][]byte
}

func NewManager(ctx context.Context, c rpc.Client, d Daemon) *Manager {
	m := &Manager{
		ctx:         ctx,
		c:           c,
		d:           d,
		syncFlags:   make(map[int64]int32),
		syncCookies: make(map[int64][]byte),
	}
	m.c.Register(service.MethodMessagePushNotify, m.handlePushNotifyRequest)
	return m
}

func (m *Manager) getSyncFlag(uin int64) int32 {
	return m.syncFlags[uin]
}

func (m *Manager) setSyncFlag(uin int64, flag int32) {
	m.syncFlags[uin] = flag
}

func (m *Manager) getSyncCookie(uin int64) []byte {
	return m.syncCookies[uin]
}

func (m *Manager) setSyncCookie(uin int64, cookie []byte) {
	m.syncCookies[uin] = cookie
}
