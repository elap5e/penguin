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

	"github.com/elap5e/penguin/daemon/auth"
	"github.com/elap5e/penguin/daemon/message/pb"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Daemon interface {
	GetAuthManager() *auth.Manager
	OnRecvMessage(uin int64, head *pb.MsgCommon_MsgHead, body *pb.IMMsgBody_MsgBody) error
}

type Manager struct {
	ctx context.Context

	c rpc.Client
	d Daemon
}

func NewManager(ctx context.Context, c rpc.Client, d Daemon) *Manager {
	m := &Manager{
		ctx: ctx,
		c:   c,
		d:   d,
	}
	m.c.Register(service.MethodServiceConfigPushDomain, m.handleConfigPushDomain)
	m.c.Register(service.MethodServiceConfigPushRequest, m.handleConfigPushRequest)
	m.c.Register(service.MethodServiceOnlinePushChatMessage, m.handleOnlinePushMessage)
	m.c.Register(service.MethodServiceOnlinePushUserMessage, m.handleOnlinePushMessage)
	m.c.Register(service.MethodServiceOnlinePushTicketExpired, m.handleOnlinePushTicketExpired)
	return m
}
