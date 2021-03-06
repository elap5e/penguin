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

package service

import (
	"context"
	"time"

	"github.com/elap5e/penguin/daemon/auth"
	"github.com/elap5e/penguin/daemon/chat"
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/fake"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Daemon interface {
	Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error
	Register(serviceMethod string, handler rpc.Handler) error

	GetAuthManager() *auth.Manager
	GetChatManager() *chat.Manager
	GetFakeSource(uin int64) *fake.Source
	OnRecvMessage(uin int64, head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) error
}

type Manager struct {
	Daemon
	ctx context.Context

	timers map[int64]*time.Timer
}

func NewManager(ctx context.Context, d Daemon) *Manager {
	m := &Manager{
		Daemon: d,
		ctx:    ctx,
		timers: make(map[int64]*time.Timer),
	}
	m.Daemon.Register(service.MethodServiceConfigPushDomain, m.handleConfigPushDomain)
	m.Daemon.Register(service.MethodServiceConfigPushRequest, m.handleConfigPushRequest)
	m.Daemon.Register(service.MethodServiceOnlinePushRequest, m.handleOnlinePushRequest)
	m.Daemon.Register(service.MethodServiceOnlinePushChatMessage, m.handleOnlinePushMessage)
	m.Daemon.Register(service.MethodServiceOnlinePushChatSerivce, m.handleOnlinePushService)
	m.Daemon.Register(service.MethodServiceOnlinePushUserMessage, m.handleOnlinePushMessage)
	m.Daemon.Register(service.MethodServiceOnlinePushTicketExpired, m.handleOnlinePushTicketExpired)
	m.Daemon.Register(service.MethodServicePushForceOffline, m.handlePushForceOffline)
	m.Daemon.Register(service.MethodServicePushLoginNotify, m.handlePushLoginNotify)
	return m
}
