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
	"encoding/base64"
	"encoding/json"

	"github.com/elap5e/penguin/daemon/channel/pb"
	"github.com/elap5e/penguin/pkg/log"
)

func (m *Manager) GetChannelRoles(uin, channelID int64) (*pb.Channel_GetChannelRolesResponse, error) {
	req := pb.Channel_GetChannelRolesRequest{
		ChannelId: channelID,
		Field3: &pb.Channel_GetChannelRolesRequest_Field3{
			Field1: 1,
			Field2: 1,
		},
	}
	resp := pb.Channel_GetChannelRolesResponse{}
	p, err := m.request(uin, 4121, 1, &req, &resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Debug("dump base64: %s", base64.RawStdEncoding.EncodeToString(p))
	p, _ = json.Marshal(&resp)
	log.Debug("dump: %s", string(p))
	return &resp, nil
}
