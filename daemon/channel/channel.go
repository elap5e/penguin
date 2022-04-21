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
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/account"
	"github.com/elap5e/penguin/daemon/channel/pb"
	"github.com/elap5e/penguin/fake"
	"github.com/elap5e/penguin/pkg/encoding/oidb"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type Daemon interface {
	Call(serviceMethod string, args *rpc.Args, reply *rpc.Reply) error
	Register(serviceMethod string, handler rpc.Handler) error

	GetFakeSource(uin int64) *fake.Source
	GetAccountManager() *account.Manager
	OnRecvChannelMessage(id int64, recv *pb.Common_Msg) error
}

type Manager struct {
	context.Context
	Daemon
	Store
}

func NewManager(ctx context.Context, daemon Daemon, store Store) *Manager {
	m := &Manager{
		Context: ctx,
		Daemon:  daemon,
		Store:   store,
	}
	m.Register(service.MethodChannelPushMessage, m.handlePushMessage)
	m.Register(service.MethodChannelPushFirstView, m.handlePushFirstView)
	return m
}

func (m *Manager) request(uin int64, cmd, svc uint32, req proto.Message, resp proto.Message) (p []byte, err error) {
	if p, err = proto.Marshal(req); err != nil {
		return nil, err
	}
	if p, err = oidb.Marshal(&oidb.Data{
		Command: cmd,
		Service: svc,
		Payload: p,
		Client:  "android " + m.GetFakeSource(uin).App.Version,
	}); err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.Call(service.MethodOidbSvcTrpcTcp(cmd, svc), &args, &reply); err != nil {
		return nil, err
	}
	data := oidb.Data{}
	if err = oidb.Unmarshal(reply.Payload, &data); err != nil {
		return nil, err
	}
	if data.Result != 0 {
		return nil, fmt.Errorf("error(%d): %s", data.Result, data.Message)
	}
	// log.Debug("dump base64: %s", base64.RawStdEncoding.EncodeToString(p))
	if err := proto.Unmarshal(data.Payload, resp); err != nil {
		return nil, err
	}
	// p, _ = json.Marshal(&resp)
	// log.Debug("dump: %s", string(p))
	return data.Payload, nil
}
