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
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/channel/pb"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (m *Manager) handlePushMessage(reply *rpc.Reply) (*rpc.Args, error) {
	push := pb.MsgPush_MsgOnlinePush{}
	if err := proto.Unmarshal(reply.Payload, &push); err != nil {
		return nil, err
	}
	for _, v := range push.GetMsgs() {
		_ = m.d.OnRecvChannelMessage(reply.Uin, v)
	}
	return nil, nil
}