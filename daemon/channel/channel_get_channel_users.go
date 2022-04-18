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

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/channel/pb"
	"github.com/elap5e/penguin/pkg/log"
)

var (
	GetChannelUsersRequestField4 = &pb.Channel_GetChannelUsersRequest_Field4{
		Field1:  1,
		Field2:  1,
		Field3:  1,
		Field4:  1,
		Field5:  1,
		Field6:  1,
		Field7:  1,
		Field8:  1,
		Field21: 1,
		Field22: 1,
		Field25: 1,
	}
)

func (m *Manager) GetChannelUsers(uin, channelID, start int64, limit int32, cookie ...[]byte) (*pb.Channel_GetChannelUsersResponse, error) {
	req := pb.Channel_GetChannelUsersRequest{
		ChannelId: channelID,
		Field2:    3,
		Field3:    0,
		Field4:    GetChannelUsersRequestField4,
		Start:     start,
		Limit:     limit,
	}
	if len(cookie) > 0 {
		req.Cookie = cookie[0]
	}
	resp := pb.Channel_GetChannelUsersResponse{}
	p, err := m.request(uin, 3931, 1, &req, &resp)
	if err != nil {
		return nil, err
	}
	log.Debug("dump base64: %s", base64.RawStdEncoding.EncodeToString(p))
	if err := m.onGetChannelUsers(&resp); err != nil {
		log.Printf("error: %v", err)
	}
	return &resp, nil
}

func (m *Manager) GetChannelUsersByIDs(uin, channelID int64, tinyIDs []int64, start int64, limit int32, cookie ...[]byte) (*pb.Channel_GetChannelUsersResponse, error) {
	req := pb.Channel_GetChannelUsersRequest{
		ChannelId: channelID,
		Field2:    3,
		Field3:    0,
		Field4:    GetChannelUsersRequestField4,
		Start:     start,
		Limit:     limit,
		TinyIds:   tinyIDs,
	}
	if len(cookie) > 0 {
		req.Cookie = cookie[0]
	}
	resp := pb.Channel_GetChannelUsersResponse{}
	p, err := m.request(uin, 3931, 1, &req, &resp)
	if err != nil {
		return nil, err
	}
	log.Debug("dump base64: %s", base64.RawStdEncoding.EncodeToString(p))
	p, _ = json.Marshal(&resp)
	log.Debug("dump: %s", string(p))
	return &resp, nil
}

func (m *Manager) GetChannelUserRoles(uin, channelID, tinyID int64) (*pb.Channel_GetChannelUserRolesResponse, error) {
	req := pb.Channel_GetChannelUserRolesRequest{
		ChannelId: channelID,
		TinyId:    tinyID,
		Field4: &pb.Channel_GetChannelUserRolesRequest_Field4{
			Field1: 1,
			Field2: 2,
			Field3: 3,
		},
	}
	resp := pb.Channel_GetChannelUserRolesResponse{}
	p, err := m.request(uin, 4119, 1, &req, &resp)
	if err != nil {
		return nil, err
	}
	log.Debug("dump base64: %s", base64.RawStdEncoding.EncodeToString(p))
	p, _ = json.Marshal(&resp)
	log.Debug("dump: %s", string(p))
	return &resp, nil
}

func (m *Manager) onGetChannelUsers(resp *pb.Channel_GetChannelUsersResponse) error {
	p, _ := json.Marshal(&resp)
	log.Debug("dump: %s", string(p))
	if err := m.setChannelUsers(resp.ChannelId, penguin.AccountTypeChannel, resp.Owner); err != nil {
		return err
	}
	if err := m.setChannelUsers(resp.ChannelId, penguin.AccountTypeChannelBot, resp.Bots...); err != nil {
		return err
	}
	if err := m.setChannelUsers(resp.ChannelId, penguin.AccountTypeChannel, resp.Users...); err != nil {
		return err
	}
	return nil
}

func (m *Manager) setChannelUsers(channelID int64, typ penguin.AccountType, users ...*pb.Channel_User) error {
	for _, v := range users {
		account := penguin.Account{
			ID:       v.GetTinyId(),
			Type:     typ,
			Username: v.GetUsername(),
		}
		_, _ = m.GetAccountManager().SetChannelAccount(account.ID, &account)
		user := penguin.User{
			Account: &account,
			Display: v.GetDisplay(),
		}
		_, _ = m.SetUser(channelID, account.ID, &user)
		p, _ := json.Marshal(user)
		log.Debug("channel:%d:user:%d:%s", channelID, account.ID, p)
	}
	return nil
}
