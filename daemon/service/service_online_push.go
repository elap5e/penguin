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
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

func (m *Manager) handleOnlinePushMessage(reply *rpc.Reply) (*rpc.Args, error) {
	if _, err := m.d.GetAuthManager().SignInChangeToken(reply.Uin); err != nil {
		return nil, err
	}
	if _, err := m.RegisterAppRegister(reply.Uin); err != nil {
		return nil, err
	}
	return &rpc.Args{
		Version:       rpc.VersionSimple,
		Uin:           reply.Uin,
		Seq:           reply.Seq,
		ServiceMethod: service.MethodServiceOnlinePushTicketExpired,
	}, nil
}

func (m *Manager) handleOnlinePushTicketExpired(reply *rpc.Reply) (*rpc.Args, error) {
	if _, err := m.d.GetAuthManager().SignInChangeToken(reply.Uin); err != nil {
		return nil, err
	}
	if _, err := m.RegisterAppRegister(reply.Uin); err != nil {
		return nil, err
	}
	return &rpc.Args{
		Version:       rpc.VersionSimple,
		Uin:           reply.Uin,
		Seq:           reply.Seq,
		ServiceMethod: service.MethodServiceOnlinePushTicketExpired,
	}, nil
}
