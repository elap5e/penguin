package chat

import (
	"encoding/json"

	"github.com/elap5e/penguin"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type GetGroupUsersRequest struct {
	Uin                 int64 `jce:"0" json:"uin"`
	GroupCode           int64 `jce:"1" json:"group_code"`
	NextUin             int64 `jce:"2" json:"next_uin"`
	GroupUin            int64 `jce:"3" json:"group_uin"`
	Version             int64 `jce:"4" json:"version"`
	ReqType             int64 `jce:"5" json:"req_type"`
	GetListAppointTime  int64 `jce:"6" json:"get_list_appoint_time"`
	RichCardNameVersion int8  `jce:"7" json:"rich_card_name_version"`
}

type GetGroupUsersResponse struct {
	Uin         int64       `jce:"0" json:",omitempty"`
	GroupCode   int64       `jce:"1" json:",omitempty"`
	GroupUin    int64       `jce:"2" json:",omitempty"`
	GroupUsers  []GroupUser `jce:"3" json:",omitempty"`
	NextUin     int64       `jce:"4" json:",omitempty"`
	Result      int32       `jce:"5" json:",omitempty"`
	ErrorCode   int16       `jce:"6" json:",omitempty"`
	OfficeMode  int64       `jce:"7" json:",omitempty"`
	NextGetTime int64       `jce:"8" json:",omitempty"`
}

type GroupUser struct {
	MemberUin              int64          `jce:"0" json:",omitempty"`
	FaceID                 int16          `jce:"1" json:",omitempty"`
	Age                    int8           `jce:"2" json:",omitempty"`
	Gender                 int8           `jce:"3" json:",omitempty"`
	Nick                   string         `jce:"4" json:",omitempty"`
	Status                 int8           `jce:"5" json:",omitempty"`
	ShowName               string         `jce:"6" json:",omitempty"`
	Name                   string         `jce:"8" json:",omitempty"`
	Gender2                int8           `jce:"9" json:",omitempty"`
	Phone                  string         `jce:"10" json:",omitempty"`
	Email                  string         `jce:"11" json:",omitempty"`
	Memo                   string         `jce:"12" json:",omitempty"`
	AutoRemark             string         `jce:"13" json:",omitempty"`
	MemberLevel            int64          `jce:"14" json:",omitempty"`
	JoinTime               int64          `jce:"15" json:",omitempty"`
	LastSpeakTime          int64          `jce:"16" json:",omitempty"`
	CreditLevel            int64          `jce:"17" json:",omitempty"`
	Flag                   int64          `jce:"18" json:",omitempty"`
	FlagExt                int64          `jce:"19" json:",omitempty"`
	Point                  int64          `jce:"20" json:",omitempty"`
	Concerned              int8           `jce:"21" json:",omitempty"`
	Shielded               int8           `jce:"22" json:",omitempty"`
	SpecialTitle           string         `jce:"23" json:",omitempty"`
	SpecialTitleExpireTime int64          `jce:"24" json:",omitempty"`
	Job                    string         `jce:"25" json:",omitempty"`
	ApolloFlag             int8           `jce:"26" json:",omitempty"`
	ApolloTimestamp        int64          `jce:"27" json:",omitempty"`
	GlobalGroupLevel       int64          `jce:"28" json:",omitempty"`
	TitleId                int64          `jce:"29" json:",omitempty"`
	ShutupTimestap         int64          `jce:"30" json:",omitempty"`
	GlobalGroupPoint       int64          `jce:"31" json:",omitempty"`
	QZoneUserInfo          *QZoneUserInfo `jce:"32" json:",omitempty"`
	RichCardNameVer        int8           `jce:"33" json:",omitempty"`
	VipType                int64          `jce:"34" json:",omitempty"`
	VipLevel               int64          `jce:"35" json:",omitempty"`
	BigClubLevel           int64          `jce:"36" json:",omitempty"`
	BigClubFlag            int64          `jce:"37" json:",omitempty"`
	Nameplate              int64          `jce:"38" json:",omitempty"`
	GroupHonor             []byte         `jce:"39" json:",omitempty"`
	Remark                 []byte         `jce:"40" json:",omitempty"`
	RichFlag               int8           `jce:"41" json:",omitempty"`
}

type QZoneUserInfo struct {
	StarState  int32             `jce:"0" json:",omitempty"`
	ExtendInfo map[string]string `jce:"1" json:",omitempty"`
}

func (m *Manager) GetGroupUsers(uin, id int64, v ...int64) (*GetGroupUsersResponse, error) {
	nextUin := int64(0)
	if len(v) != 0 {
		nextUin = v[0]
	}
	return m.requestGetGroupUsers(uin, &GetGroupUsersRequest{
		Uin:                 uin,
		GroupCode:           id,
		NextUin:             nextUin,
		GroupUin:            id,
		Version:             3,
		ReqType:             0,
		GetListAppointTime:  0,
		RichCardNameVersion: 1,
	})
}

func (m *Manager) requestGetGroupUsers(uin int64, req *GetGroupUsersRequest) (*GetGroupUsersResponse, error) {
	p, err := uni.Marshal(&uni.Data{
		Version:     3,
		ServantName: "mqq.IMService.FriendListServiceServantObj",
		FuncName:    "GetTroopMemberListReq",
	}, map[string]interface{}{
		"GTML": req,
	})
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{
		Version: rpc.VersionSimple,
		Uin:     uin,
		Payload: p,
	}, rpc.Reply{}
	if err = m.c.Call(service.MethodChatGetGroupUsers, &args, &reply); err != nil {
		return nil, err
	}
	data, resp := uni.Data{}, GetGroupUsersResponse{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"GTMLRESP": &resp,
	}); err != nil {
		return nil, err
	}
	for _, v := range resp.GroupUsers {
		account := penguin.Account{
			ID:       v.MemberUin,
			Type:     penguin.AccountTypeDefault,
			Username: v.Nick,
		}
		_, _ = m.d.GetAccountManager().SetAccount(account.ID, &account)
		user := penguin.User{
			Account: &account,
			Display: string(v.Remark),
		}
		_, _ = m.SetUser(resp.GroupCode, account.ID, &user)
		p, _ := json.Marshal(user)
		log.Debug("chat:%d:user:%d:%s", resp.GroupCode, account.ID, p)
	}
	return &resp, nil
}
