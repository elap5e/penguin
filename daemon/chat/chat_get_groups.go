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

package chat

import (
	"encoding/json"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type GetGroupsRequest struct {
	Uin               int64       `jce:"0" json:"uin"`
	GetMSFMessageFlag int8        `jce:"1" json:"get_msf_message_flag"`
	Cookie            []byte      `jce:"2" json:"cookie"`
	Groups            []*GroupSeq `jce:"3" json:"groups"`
	GroupFlagExtra    int8        `jce:"4" json:"group_flag_extra"`
	Version           int32       `jce:"5" json:"version"`
	CompanyID         int64       `jce:"6" json:"company_id"`
	VersionNumber     int64       `jce:"7" json:"version_number"`
	GetLongGroupName  int8        `jce:"8" json:"get_long_group_name"`
}

type GroupSeq struct {
	Code         int64 `jce:"0" json:"code"`
	InfoSeq      int64 `jce:"1" json:"info_seq"`
	FlagExtra    int64 `jce:"2" json:"flag_extra"`
	RankSeq      int64 `jce:"3" json:"rank_seq"`
	InfoExtraSeq int64 `jce:"4" json:"info_extra_seq"`
}

type GetGroupsResponse struct {
	Uin             int64            `jce:"0" json:",omitempty"`
	GroupCount      int16            `jce:"1" json:",omitempty"`
	Result          int32            `jce:"2" json:",omitempty"`
	ErrorCode       int16            `jce:"3" json:",omitempty"`
	Cookie          []byte           `jce:"4" json:",omitempty"`
	Groups          []*Group         `jce:"5" json:",omitempty"`
	GroupsDelete    []*Group         `jce:"6" json:",omitempty"`
	GroupRanks      []*GroupRank     `jce:"7" json:",omitempty"`
	FavouriteGroups []*FavoriteGroup `jce:"8" json:",omitempty"`
	GroupsExtra     []*Group         `jce:"9" json:",omitempty"`
	GroupExtra      []int64          `jce:"10" json:",omitempty"`
}

type FavoriteGroup struct {
	GroupCode     int64 `jce:"0" json:",omitempty"`
	Timestamp     int64 `jce:"1" json:",omitempty"`
	SNSFlag       int64 `jce:"2" json:",omitempty"`
	OpenTimestamp int64 `jce:"3" json:",omitempty"`
}

type GroupRank struct {
	GroupCode            int64        `jce:"0" json:",omitempty"`
	GroupRankSysFlag     int8         `jce:"1" json:",omitempty"`
	GroupRankUserFlag    int8         `jce:"2" json:",omitempty"`
	RankMap              []*LevelRank `jce:"3" json:",omitempty"`
	GroupRankSeq         int64        `jce:"4" json:",omitempty"`
	OwnerName            string       `jce:"5" json:",omitempty"`
	AdminName            string       `jce:"6" json:",omitempty"`
	OfficeMode           int64        `jce:"7" json:",omitempty"`
	GroupRankUserFlagNew int8         `jce:"8" json:",omitempty"`
	RankMapNew           []*LevelRank `jce:"9" json:",omitempty"`
}

type LevelRank struct {
	Level int64  `jce:"0" json:",omitempty"`
	Rank  string `jce:"1" json:",omitempty"`
}

type Group struct {
	Uin                   int64  `jce:"0" json:",omitempty"`
	Code                  int64  `jce:"1" json:",omitempty"`
	Flag                  int8   `jce:"2" json:",omitempty"`
	InfoSeq               int64  `jce:"3" json:",omitempty"`
	Name                  string `jce:"4" json:",omitempty"`
	Memo                  string `jce:"5" json:",omitempty"`
	FlagExt               int64  `jce:"6" json:",omitempty"`
	RankSeq               int64  `jce:"7" json:",omitempty"`
	CertificationType     int64  `jce:"8" json:",omitempty"`
	ShutupTimestamp       int64  `jce:"9" json:",omitempty"`
	MyShutupTimestamp     int64  `jce:"10" json:",omitempty"`
	CmdUinUinFlag         int64  `jce:"11" json:",omitempty"`
	AdditionalFlag        int64  `jce:"12" json:",omitempty"`
	TypeFlag              int64  `jce:"13" json:",omitempty"`
	SecType               int64  `jce:"14" json:",omitempty"`
	SecTypeInfo           int64  `jce:"15" json:",omitempty"`
	ClassExt              int64  `jce:"16" json:",omitempty"`
	AppPrivilegeFlag      int64  `jce:"17" json:",omitempty"`
	SubscriptionUin       int64  `jce:"18" json:",omitempty"`
	MemberNum             int64  `jce:"19" json:",omitempty"`
	MemberNumSeq          int64  `jce:"20" json:",omitempty"`
	MemberCardSeq         int64  `jce:"21" json:",omitempty"`
	FlagExt3              int64  `jce:"22" json:",omitempty"`
	OwnerUin              int64  `jce:"23" json:",omitempty"`
	IsConfGroup           bool   `jce:"24" json:",omitempty"`
	IsModifyConfGroupFace bool   `jce:"25" json:",omitempty"`
	IsModifyConfGroupName bool   `jce:"26" json:",omitempty"`
	CmduinJoinTime        int64  `jce:"27" json:",omitempty"`
	CompanyID             int64  `jce:"28" json:",omitempty"`
	MaxMemberNum          int64  `jce:"29" json:",omitempty"`
	CmdUinMask            int64  `jce:"30" json:",omitempty"`
	HLGuildAppid          int64  `jce:"31" json:",omitempty"`
	HLGuildSubType        int64  `jce:"32" json:",omitempty"`
	CmdUinRingtoneID      int64  `jce:"33" json:",omitempty"`
	CmdUinFlagEx2         int64  `jce:"34" json:",omitempty"`
	FlagExt4              int64  `jce:"35" json:",omitempty"`
	AppealDeadline        int64  `jce:"36" json:",omitempty"`
	Flag1                 int64  `jce:"37" json:",omitempty"`
	Remark                []byte `jce:"38" json:",omitempty"`
}

func (m *Manager) GetGroups(uin int64) (*GetGroupsResponse, error) {
	cookie, _ := m.GetCookie(uin)
	return m.requestGetGroups(uin, &GetGroupsRequest{
		Uin:               uin,
		GetMSFMessageFlag: 0,
		Cookie:            cookie,
		Groups:            []*GroupSeq{},
		GroupFlagExtra:    1,
		Version:           9,
		CompanyID:         0,
		VersionNumber:     1,
		GetLongGroupName:  1,
	})
}

func (m *Manager) requestGetGroups(uin int64, req *GetGroupsRequest) (*GetGroupsResponse, error) {
	p, err := uni.Marshal(&uni.Data{
		Version:     3,
		ServantName: "mqq.IMService.FriendListServiceServantObj",
		FuncName:    "GetTroopListReqV2Simplify",
	}, map[string]any{
		"GetTroopListReqV2Simplify": req,
	})
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.Call(service.MethodChatGetGroups, &args, &reply); err != nil {
		return nil, err
	}
	data, resp := uni.Data{}, GetGroupsResponse{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"GetTroopListRespV2": &resp,
	}); err != nil {
		return nil, err
	}
	for _, v := range resp.Groups {
		chat := penguin.Chat{
			ID:      v.Code,
			Type:    penguin.ChatTypeGroup,
			Title:   v.Name,
			Display: string(v.Remark),
		}
		_, _ = m.SetChat(chat.ID, &chat)
		p, _ := json.Marshal(chat)
		log.Debug("chat:%d:%s", chat.ID, p)
		_ = m.GetGroupUsersAll(uin, chat.ID)
	}
	return &resp, nil
}
