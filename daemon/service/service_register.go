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
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/constant"
	"github.com/elap5e/penguin/daemon/service/pb"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type StatusType uint32

var (
	StatusTypeOnline     StatusType = 11
	StatusTypeOffline    StatusType = 21
	StatusTypeAway       StatusType = 31
	StatusTypeInvisiable StatusType = 41
	StatusTypeBusy       StatusType = 50
	StatusTypeQMe        StatusType = 60
	StatusTypeDND        StatusType = 70 // do not disturb
	StatusTypeReceiveMsg StatusType = 95 // receive offline message
)

type RegisterType uint8

const (
	RegisterTypeMSFBoot RegisterType = iota
	RegisterTypeAppRegister
	RegisterTypeSetOnlineStatus
	RegisterTypeMSFHeartbeatTimeout
	RegisterTypeMSFNetworkChange
	RegisterTypeServerPush
	RegisterTypeFillRegProxy
	RegisterTypeCreateDefaultRegInfo
)

type RegisterRequest struct {
	Uin          int64      `jce:"0" json:"uin"`
	Bid          int64      `jce:"1" json:"bid"`
	ConnectType  bool       `jce:"2" json:"connect_type"`
	Other        string     `jce:"3" json:"other"`
	Status       StatusType `jce:"4" json:"status"`
	OnlinePush   bool       `jce:"5" json:"online_push"`
	IsOnline     bool       `jce:"6" json:"is_online"`
	IsShowOnline bool       `jce:"7" json:"is_show_online"`
	KickPC       bool       `jce:"8" json:"kick_pc"`
	KickWeak     bool       `jce:"9" json:"kick_weak"`
	Timestamp    int64      `jce:"10" json:"timestamp"`
	SDKVersion   int64      `jce:"11" json:"sdk_version"`
	NetworkType  bool       `jce:"12" json:"network_type"`
	BuildVersion string     `jce:"13" json:"build_version"`
	RegisterType bool       `jce:"14" json:"register_type"`
	DeviceParam  []byte     `jce:"15" json:"device_param"`
	GUID         []byte     `jce:"16" json:"guid"`
	LocaleID     uint32     `jce:"17" json:"locale_id"`
	SlientPush   bool       `jce:"18" json:"slient_push"`
	DeviceName   string     `jce:"19" json:"device_name"`
	DeviceType   string     `jce:"20" json:"device_type"`
	OSVersion    string     `jce:"21" json:"os_version"`
	OpenPush     bool       `jce:"22" json:"open_push"`
	LargeSeq     int64      `jce:"23" json:"large_seq"`

	LastWatchStart int64           `jce:"24" json:"last_watch_start"`
	BindUinList    []*BindUinInfo  `jce:"25" json:"bind_uin_list"`
	OldSSOIP       int64           `jce:"26" json:"old_sso_ip"`
	NewSSOIP       int64           `jce:"27" json:"new_sso_ip"`
	ChannelID      string          `jce:"28" json:"channel_id"`
	CPID           int64           `jce:"29" json:"cpid"`
	VendorName     string          `jce:"30" json:"vendor_name"`
	VendorOSName   string          `jce:"31" json:"vendor_os_name"`
	IOSIDFA        string          `jce:"32" json:"ios_idfa"`
	Body0x769      []byte          `jce:"33" json:"body_0x796"`
	IsSetStatus    bool            `jce:"34" json:"is_set_status"`
	ServerPayload  []byte          `jce:"35" json:"server_payload"`
	SetMute        bool            `jce:"36" json:"set_mute"`
	NotifySwitch   bool            `jce:"37" json:"notify_switch"`
	ExtraStatus    int64           `jce:"38" json:"extra_status"`
	BatteryStatus  uint32          `jce:"39" json:"battery_status"`
	TimActiveFlag  bool            `jce:"40" json:"tim_active_flag"`
	BindUinNotify  bool            `jce:"41" json:"bind_uin_notify"`
	VendorPushInfo *VendorPushInfo `jce:"42" json:"vendor_push_info,omitempty"`
	VendorDeviceID int64           `jce:"43" json:"vendor_device_id"`
	CustomStatus   []byte          `jce:"45" json:"custom_status"`
}

type BindUinInfo struct {
	Uin          int64      `jce:"0" json:"uin"`
	CustomStatus []byte     `jce:"1" json:"custom_status"`
	Status       StatusType `jce:"2" json:"status"`
}

type VendorPushInfo struct {
	Type uint64 `jce:"0" json:"type"`
}

type RegisterResponse struct {
	Uin            int64  `jce:"0" json:"uin"`
	Bid            int64  `jce:"1" json:"bid"`
	ReplyCode      uint8  `jce:"2" json:"reply_code"`
	Result         string `jce:"3" json:"result"`
	ServerTime     int64  `jce:"4" json:"server_time"`
	LogQQ          bool   `jce:"5" json:"log_qq"`
	NeedKick       bool   `jce:"6" json:"need_kick"`
	UpdateFlag     bool   `jce:"7" json:"update_flag"`
	Timestamp      int64  `jce:"8" json:"timestamp"`
	CrashFlag      bool   `jce:"9" json:"crash_flag"`
	ClientIP       string `jce:"10" json:"client_ip"`
	ClientPort     int32  `jce:"11" json:"client_port"`
	HelloInterval  int32  `jce:"12" json:"hello_interval"`
	LargeSeq       int32  `jce:"13" json:"large_seq"`
	LargeSeqUpdate bool   `jce:"14" json:"large_seq_update"`

	Body0x769                []byte `jce:"15" json:"body_0x796"`
	Status                   int32  `jce:"16" json:"status"`
	ExtraStatus              int64  `jce:"17" json:"extra_status"`
	ClientBatteryGetInterval int64  `jce:"18" json:"client_battery_get_interval"`
	ClientAutoStatusInterval int64  `jce:"19" json:"client_auto_status_interval"`
	CustomStatus             []byte `jce:"21" json:"custom_status"`
}

type RegisterPush struct {
	Uin           int64      `jce:"1" json:"uin"`
	PushIDs       []int64    `jce:"2" json:"push_ids"`
	Status        StatusType `jce:"3" json:"status"`
	KickPC        bool       `jce:"4" json:"kick_pc"`
	KickWeak      bool       `jce:"5" json:"kick_weak"`
	Timestamp     int64      `jce:"6" json:"timestamp"`
	LargeSeq      int64      `jce:"7" json:"large_seq"`
	ExtraStatus   int64      `jce:"8" json:"extra_status"`
	BatteryCap    int32      `jce:"9" json:"battery_cap"`
	PowerConnect  int32      `jce:"10" json:"power_connect"`
	BindUinNotify bool       `jce:"11" json:"bind_uin_notify"`
}

// RegPushReason.appRegister
func (m *Manager) RegisterAppRegister(uin int64) (*RegisterResponse, error) {
	return m.Register(uin, StatusTypeOnline, false, RegisterTypeAppRegister)
}

// RegPushReason.setOnlineStatus
func (m *Manager) RegisterSetOnlineStatus(uin int64, status StatusType, client ...bool) (*RegisterResponse, error) {
	return m.Register(uin, status, false, RegisterTypeSetOnlineStatus, client...)
}

func (m *Manager) Register(uin int64, status StatusType, kick bool, typ RegisterType, client ...bool) (*RegisterResponse, error) {
	return m.requestRegisterPush(&RegisterPush{
		Uin:           uin,
		PushIDs:       []int64{1, 2, 4}, // constant
		Status:        status,
		KickPC:        kick,
		KickWeak:      false, // constant
		Timestamp:     0,     // service_register_time
		LargeSeq:      0,     // friend_list_seq
		ExtraStatus:   0,     // constant
		BatteryCap:    0,     // constant
		PowerConnect:  -1,    // constant
		BindUinNotify: false, // sub_account_notify
	}, typ, client...)
}

func (m *Manager) requestRegisterPush(push *RegisterPush, typ RegisterType, client ...bool) (*RegisterResponse, error) {
	fake := m.c.GetFakeSource(push.Uin)
	bid := int64(0)
	for _, id := range push.PushIDs {
		bid |= id
	}
	body, err := proto.Marshal(&pb.Oidb_0X769_ReqBody{
		ConfigList: []*pb.Oidb_0X769_ConfigSeq{
			{Type: 46, Version: 0}, // key_config_version_patch
			{Type: 283, Version: 0},
		},
	})
	if err != nil {
		return nil, err
	}
	return m.requestRegister(&RegisterRequest{
		Uin:            push.Uin,
		Bid:            bid,
		ConnectType:    false, // constant
		Other:          "",    // constant
		Status:         push.Status,
		OnlinePush:     false, // constant
		IsOnline:       false, // constant
		IsShowOnline:   false, // constant
		KickPC:         push.KickPC,
		KickWeak:       push.KickWeak,
		Timestamp:      push.Timestamp,
		SDKVersion:     int64(fake.Device.OS.SDKVersion),
		NetworkType:    true, // false:mobile, ture:wifi
		BuildVersion:   "",   // constant
		RegisterType:   !(typ == RegisterTypeAppRegister || typ == RegisterTypeSetOnlineStatus || typ == RegisterTypeFillRegProxy || typ == RegisterTypeCreateDefaultRegInfo),
		DeviceParam:    nil, // constant
		GUID:           fake.Device.GUID[:],
		LocaleID:       constant.LocaleID, // constant
		SlientPush:     false,             // constant
		DeviceName:     fake.Device.OS.BuildModel,
		DeviceType:     fake.Device.OS.BuildModel,
		OSVersion:      fake.Device.OS.Version,
		OpenPush:       true, // constant
		LargeSeq:       push.LargeSeq,
		LastWatchStart: 0,                // constant
		BindUinList:    []*BindUinInfo{}, // constant
		OldSSOIP:       0,                // constant
		NewSSOIP:       -1,               // constant
		ChannelID:      "",               // constant
		CPID:           0,                // constant
		VendorName:     "MIUI",
		VendorOSName:   "V13",
		IOSIDFA:        "", // constant
		Body0x769:      body,
		IsSetStatus:    typ == RegisterTypeSetOnlineStatus,
		ServerPayload:  nil,
		SetMute:        false, // qqsetting_qrlogin_set_mute
		NotifySwitch:   true,  // constant
		ExtraStatus:    push.ExtraStatus,
		BatteryStatus:  0,     // constant
		TimActiveFlag:  false, // constant
		BindUinNotify:  push.BindUinNotify,
		VendorPushInfo: &VendorPushInfo{}, // constant
		VendorDeviceID: 0,                 // constant
		CustomStatus:   nil,               // constant
	}, client...)
}

func (m *Manager) requestRegister(req *RegisterRequest, client ...bool) (*RegisterResponse, error) {
	p, err := uni.Marshal(&uni.Data{
		Version:     3,
		ServantName: "PushService",
		FuncName:    "SvcReqRegister",
	}, map[string]any{
		"SvcReqRegister": req,
	})
	if err != nil {
		return nil, err
	}
	args, reply := rpc.Args{Uin: req.Uin, Payload: p}, rpc.Reply{}
	serviceMethod := service.MethodServiceRegister
	if len(client) != 0 && client[0] {
		serviceMethod = service.MethodServiceSetStatusFromClient
	}
	if err = m.c.Call(serviceMethod, &args, &reply); err != nil {
		return nil, err
	}
	data, resp := uni.Data{}, RegisterResponse{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"SvcRespRegister": &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}
