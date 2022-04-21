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
	"encoding/json"

	"github.com/elap5e/penguin/daemon/message/dto"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func (m *Manager) handlePushReadedRequest(reply *rpc.Reply) (*rpc.Args, error) {
	data, push := uni.Data{}, dto.PushReadedRequest{}
	if err := uni.Unmarshal(reply.Payload[4:], &data, map[string]any{
		"req": &push,
	}); err != nil {
		return nil, err
	}
	p, _ := json.Marshal(push)
	log.Debug("message.handlePushReadedRequest: %s", p)
	return nil, nil
}
