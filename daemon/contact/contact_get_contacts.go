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

package contact

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/daemon/contact/pb"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type GetContactsRequest struct {
	RequestType     int32   `jce:"0" json:"request_type"`
	IsReflush       bool    `jce:"1" json:"is_reflush"`
	Uin             int64   `jce:"2" json:"uin"`
	StartIndex      int16   `jce:"3" json:"start_index"`
	GetFriendCount  int16   `jce:"4" json:"get_friend_count"`
	GroupID         int8    `jce:"5" json:"group_id"`
	IsGetGroup      bool    `jce:"6" json:"is_get_group_info"`
	GroupStartIndex int8    `jce:"7" json:"group_start_index"`
	GetGroupCount   int8    `jce:"8" json:"get_group_count"`
	IsGetMSFGroup   bool    `jce:"9" json:"is_get_msf_group"`
	IsShowTermType  bool    `jce:"10" json:"is_show_term_type"`
	Version         int64   `jce:"11" json:"version"`
	UinList         []int64 `jce:"12" json:"uin_list"`
	AppType         int32   `jce:"13" json:"app_type"`
	IsGetDOVID      bool    `jce:"14" json:"is_get_dovid"`
	IsGetBothFlag   bool    `jce:"15" json:"is_get_both_flag"`
	OIDB0XD50       []byte  `jce:"16" json:"oidb0xd50"`
	OIDB0XD6B       []byte  `jce:"17" json:"oidb0xd6b"`
	SNSTypeList     []int64 `jce:"18" json:"sns_type_list"`
}

type GetContactsResponse struct {
	RequestType             int32                  `jce:"0" json:",omitempty"`
	IsReflush               bool                   `jce:"1" json:",omitempty"`
	Uin                     int64                  `jce:"2" json:",omitempty"`
	StartIndex              int16                  `jce:"3" json:",omitempty"`
	GetFriendCount          int16                  `jce:"4" json:",omitempty"`
	TotoalFriendCount       int16                  `jce:"5" json:",omitempty"`
	FriendCount             int16                  `jce:"6" json:",omitempty"`
	FriendInfoList          []*Contact             `jce:"7" json:",omitempty"`
	GroupID                 int8                   `jce:"8" json:",omitempty"`
	IsGetGroup              bool                   `jce:"9" json:",omitempty"`
	GroupStartIndex         int8                   `jce:"10" json:",omitempty"`
	GetGroupCount           int8                   `jce:"11" json:",omitempty"`
	TotoalGroupCount        int16                  `jce:"12" json:",omitempty"`
	GroupCount              int8                   `jce:"13" json:",omitempty"`
	Groups                  []*ContactGroup        `jce:"14" json:",omitempty"`
	Result                  int32                  `jce:"15" json:",omitempty"`
	ErrorCode               int16                  `jce:"16" json:",omitempty"`
	OnlineFriendCount       int16                  `jce:"17" json:",omitempty"`
	ServerTime              int64                  `jce:"18" json:",omitempty"`
	QQOnlineCount           int16                  `jce:"19" json:",omitempty"`
	Groups2                 []*ContactGroup        `jce:"20" json:",omitempty"`
	RespType                int8                   `jce:"21" json:",omitempty"`
	HasOtherRespFlag        int8                   `jce:"22" json:",omitempty"`
	FriendInfo              *Contact               `jce:"23" json:",omitempty"`
	ShowPcIcon              int8                   `jce:"24" json:",omitempty"`
	GetExtraSNSResponseCode int16                  `jce:"25" json:",omitempty"`
	SubServerResponseCode   *SubServerResponseCode `jce:"26" json:",omitempty"`
}

type Contact struct {
	FriendUin             int64        `jce:"0" json:",omitempty"`
	GroupID               int8         `jce:"1" json:",omitempty"`
	FaceID                int16        `jce:"2" json:",omitempty"`
	Remark                string       `jce:"3" json:",omitempty"`
	QQType                int8         `jce:"4" json:",omitempty"`
	Status                int8         `jce:"5" json:",omitempty"`
	MemberLevel           int8         `jce:"6" json:",omitempty"`
	IsMobileQQOnLine      bool         `jce:"7" json:",omitempty"`
	QQOnLineState         int8         `jce:"8" json:",omitempty"`
	IsIphoneOnline        bool         `jce:"9" json:",omitempty"`
	DetalStatusFlag       int8         `jce:"10" json:",omitempty"`
	QQOnLineStateV2       int8         `jce:"11" json:",omitempty"`
	ShowName              string       `jce:"12" json:",omitempty"`
	IsRemark              bool         `jce:"13" json:",omitempty"`
	Nick                  string       `jce:"14" json:",omitempty"`
	SpecialFlag           int8         `jce:"15" json:",omitempty"`
	IMGroupID             []byte       `jce:"16" json:",omitempty"`
	MSFGroupID            []byte       `jce:"17" json:",omitempty"`
	TermType              int32        `jce:"18" json:",omitempty"`
	VIPBaseInfo           *VIPBaseInfo `jce:"19" json:",omitempty"`
	Network               int8         `jce:"20" json:",omitempty"`
	Ring                  []byte       `jce:"21" json:",omitempty"`
	AbiFlag               int64        `jce:"22" json:",omitempty"`
	FaceAddonId           int64        `jce:"23" json:",omitempty"`
	NetworkType           int32        `jce:"24" json:",omitempty"`
	VIPFont               int64        `jce:"25" json:",omitempty"`
	IconType              int32        `jce:"26" json:",omitempty"`
	TermDesc              string       `jce:"27" json:",omitempty"`
	ColorRing             int64        `jce:"28" json:",omitempty"`
	ApolloFlag            int8         `jce:"29" json:",omitempty"`
	ApolloTimestamp       int64        `jce:"30" json:",omitempty"`
	Gender                int8         `jce:"31" json:",omitempty"`
	FounderFont           int64        `jce:"32" json:",omitempty"`
	EimId                 string       `jce:"33" json:",omitempty"`
	EimMobile             string       `jce:"34" json:",omitempty"`
	OlympicTorch          int8         `jce:"35" json:",omitempty"`
	ApolloSignTime        int64        `jce:"36" json:",omitempty"`
	LaviUin               int64        `jce:"37" json:",omitempty"`
	TagUpdateTime         int64        `jce:"38" json:",omitempty"`
	GameLastLoginTime     int64        `jce:"39" json:",omitempty"`
	GameAppID             int64        `jce:"40" json:",omitempty"`
	CardID                []byte       `jce:"41" json:",omitempty"`
	BitSet                int64        `jce:"42" json:",omitempty"`
	KingOfGloryFlag       int8         `jce:"43" json:",omitempty"`
	KingOfGloryRank       int64        `jce:"44" json:",omitempty"`
	MasterUin             string       `jce:"45" json:",omitempty"`
	LastMedalUpdateTime   int64        `jce:"46" json:",omitempty"`
	FaceStoreId           int64        `jce:"47" json:",omitempty"`
	FontEffect            int64        `jce:"48" json:",omitempty"`
	DOVID                 string       `jce:"49" json:",omitempty"`
	BothFlag              int64        `jce:"50" json:",omitempty"`
	CentiShow3DFlag       int8         `jce:"51" json:",omitempty"`
	IntimateInfo          []byte       `jce:"52" json:",omitempty"`
	ShowNameplate         int8         `jce:"53" json:",omitempty"`
	NewLoverDiamondFlag   int8         `jce:"54" json:",omitempty"`
	ExtSnsFrdData         []byte       `jce:"55" json:",omitempty"`
	MutualMarkData        []byte       `jce:"56" json:",omitempty"`
	ExtOnlineStatus       int64        `jce:"57" json:",omitempty"`
	BatteryStatus         int32        `jce:"58" json:",omitempty"`
	MusicInfo             []byte       `jce:"59" json:",omitempty"`
	PoiInfo               []byte       `jce:"60" json:",omitempty"`
	ExtOnlineBusinessInfo []byte       `jce:"61" json:",omitempty"`
}

type VIPBaseInfo struct {
	OpenInfoMap       map[uint64]*VIPOpenInfo `jce:"0" json:",omitempty"`
	NameplateVIPType  int32                   `jce:"1" json:",omitempty"`
	GrayNameplateFlag int32                   `jce:"2" json:",omitempty"`
	ExtendNameplateId string                  `jce:"3" json:",omitempty"`
}

type VIPOpenInfo struct {
	Open        bool  `jce:"0" json:",omitempty"`
	VIPType     int32 `jce:"1" json:",omitempty"`
	VIPLevel    int32 `jce:"2" json:",omitempty"`
	VIPFlag     int32 `jce:"3" json:",omitempty"`
	NameplateID int64 `jce:"4" json:",omitempty"`
}

type SubServerResponseCode struct {
	GetMutualMarkCode   int16 `jce:"0" json:",omitempty"`
	GetIntimateInfoCode int16 `jce:"1" json:",omitempty"`
}

type ContactGroup struct {
	GroupID    int8   `jce:"0" json:"group_id"`
	GroupTitle string `jce:"1" json:"group_title"`
}

func (m *Manager) GetContacts(uin int64, idx, cnt int16, gidx, gcnt int8) (*GetContactsResponse, error) {
	body0xd58, _ := proto.Marshal(&pb.OIDB0XD50_ReqBody{
		Appid:                   10002,
		ReqMusicSwitch:          1,
		ReqKsingSwitch:          1,
		ReqMutualmarkLbsshare:   1,
		ReqMutualmarkAlienation: 1,
		ReqAioQuickApp:          1,
	})
	body0xd6b, _ := proto.Marshal(&pb.OIDB0XD6B_ReqBody{})
	return m.requestGetContacts(uin, &GetContactsRequest{
		RequestType:     3,
		IsReflush:       idx == 0,
		Uin:             uin,
		StartIndex:      idx,
		GetFriendCount:  cnt,
		GroupID:         0,
		IsGetGroup:      gcnt > 0,
		GroupStartIndex: gidx,
		GetGroupCount:   gcnt,
		IsGetMSFGroup:   false,
		IsShowTermType:  true,
		Version:         32,
		UinList:         nil,
		AppType:         0,
		IsGetDOVID:      false,
		IsGetBothFlag:   false,
		OIDB0XD50:       body0xd58,
		OIDB0XD6B:       body0xd6b,
		SNSTypeList:     []int64{13580, 13581, 13582},
	})
}

func (m *Manager) requestGetContacts(uin int64, req *GetContactsRequest) (*GetContactsResponse, error) {
	p, err := uni.Marshal(&uni.Data{
		Version:     3,
		ServantName: "mqq.IMService.FriendListServiceServantObj",
		FuncName:    "GetFriendListReq",
	}, map[string]any{
		"FL": req,
	})
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.Call(service.MethodContactGetContacts, &args, &reply); err != nil {
		return nil, err
	}
	data, resp := uni.Data{}, GetContactsResponse{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"FLRESP": &resp,
	}); err != nil {
		return nil, err
	}
	for _, v := range resp.FriendInfoList {
		account := penguin.Account{
			ID:       v.FriendUin,
			Type:     penguin.AccountTypeDefault,
			Username: v.Nick,
		}
		_, _ = m.GetAccountManager().SetDefaultAccount(account.ID, &account)
		contact := penguin.Contact{
			Account: &account,
			Display: v.Remark,
		}
		_, _ = m.SetContact(uin, account.ID, &contact)
		p, _ := json.Marshal(contact)
		log.Debug("contact:%d:%s", account.ID, p)
	}
	return &resp, nil
}
