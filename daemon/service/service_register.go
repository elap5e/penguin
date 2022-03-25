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
	"github.com/elap5e/penguin/pkg/encoding/uni"
)

type RegisterRequest struct {
	Uin          int64  `jce:"0" json:"uin,omitempty"`
	Bid          int64  `jce:"1" json:"bid,omitempty"`
	ConnectType  uint8  `jce:"2" json:"connect_type,omitempty"` // constant 0x00
	Other        string `jce:"3" json:"other,omitempty"`        // constant ""
	Status       uint32 `jce:"4" json:"status,omitempty"`
	OnlinePush   bool   `jce:"5" json:"online_push,omitempty"`    // constant false
	IsOnline     bool   `jce:"6" json:"is_online,omitempty"`      // constant false
	IsShowOnline bool   `jce:"7" json:"is_show_online,omitempty"` // constant false
	KickPC       bool   `jce:"8" json:"kick_pc,omitempty"`
	KickWeak     bool   `jce:"9" json:"kick_weak,omitempty"` // constant false
	Timestamp    uint64 `jce:"10" json:"timestamp,omitempty"`
	SDKVersion   uint32 `jce:"11" json:"sdk_version,omitempty"`
	NetworkType  uint8  `jce:"12" json:"network_type,omitempty"`  // 0x00: mobile; 0x01: wifi
	BuildVersion string `jce:"13" json:"build_version,omitempty"` // constant ""
	RegisterType bool   `jce:"14" json:"register_type,omitempty"` // false: appRegister, fillRegProxy, createDefaultRegInfo; true: others
	DeviceParam  []byte `jce:"15" json:"device_param,omitempty"`  // constant nil
	GUID         []byte `jce:"16" json:"guid,omitempty"`          // placeholder
	LocaleID     uint32 `jce:"17" json:"locale_id,omitempty"`     // constant 0x00000804
	SlientPush   bool   `jce:"18" json:"slient_push,omitempty"`   // constant false
	DeviceName   string `jce:"19" json:"device_name,omitempty"`
	DeviceType   string `jce:"20" json:"device_type,omitempty"`
	OSVersion    string `jce:"21" json:"os_version,omitempty"`
	OpenPush     bool   `jce:"22" json:"open_push,omitempty"` // constant true
	LargeSeq     uint32 `jce:"23" json:"large_seq,omitempty"` // constant 0x00000000

	LastWatchStart uint32          `jce:"24" json:"last_watch_start,omitempty"`
	BindUin        []uint64        `jce:"25" json:"bind_uin,omitempty"`
	OldSSOIP       uint64          `jce:"26" json:"old_sso_ip,omitempty"`
	NewSSOIP       uint64          `jce:"27" json:"new_sso_ip,omitempty"`
	ChannelID      string          `jce:"28" json:"channel_id,omitempty"`
	CPID           uint64          `jce:"29" json:"cpid,omitempty"`
	VendorName     string          `jce:"30" json:"vendor_name,omitempty"`
	VendorOSName   string          `jce:"31" json:"vendor_os_name,omitempty"`
	IOSIDFA        string          `jce:"32" json:"ios_idfa,omitempty"`
	Body0x769      []byte          `jce:"33" json:"body_0x796,omitempty"`
	IsSetStatus    bool            `jce:"34" json:"is_set_status,omitempty"`
	ServerBuffer   []byte          `jce:"35" json:"server_buffer,omitempty"`
	SetMute        bool            `jce:"36" json:"set_mute,omitempty"`
	NotifySwitch   uint8           `jce:"37" json:"notify_switch,omitempty"`
	ExtraStatus    uint64          `jce:"38" json:"extra_status,omitempty"`
	BatteryStatus  uint32          `jce:"39" json:"battery_status,omitempty"`
	TimActiveFlag  bool            `jce:"40" json:"tim_active_flag,omitempty"`
	BindUinNotify  string          `jce:"41" json:"bind_uin_notify,omitempty"`
	VendorPushInfo *VendorPushInfo `jce:"42" json:"vendor_push_info,omitempty"`
	VendorDeviceID int64           `jce:"43" json:"vendor_device_id,omitempty"`
	CustomStatus   []byte          `jce:"45" json:"custom_status,omitempty"`
}

type VendorPushInfo struct {
	Type uint64 `jce:"0" json:"type,omitempty"`
}

type AccountSetStatusResponse struct {
	Uin            int64  `jce:"0" json:"uin,omitempty"`
	Bid            int64  `jce:"1" json:"bid,omitempty"`
	ReplyCode      uint8  `jce:"2" json:"reply_code,omitempty"`
	Result         string `jce:"3" json:"result,omitempty"`
	ServerTime     int64  `jce:"4" json:"server_time,omitempty"`
	LogQQ          bool   `jce:"5" json:"log_qq,omitempty"`
	NeedKick       bool   `jce:"6" json:"need_kick,omitempty"`
	UpdateFlag     bool   `jce:"7" json:"update_flag,omitempty"`
	Timestamp      int64  `jce:"8" json:"timestamp,omitempty"`
	CrashFlag      bool   `jce:"9" json:"crash_flag,omitempty"`
	ClientIP       string `jce:"10" json:"client_ip,omitempty"`
	ClientPort     int32  `jce:"11" json:"client_port,omitempty"`
	HelloInterval  int32  `jce:"12" json:"hello_interval,omitempty"`
	LargeSeq       int32  `jce:"13" json:"large_seq,omitempty"`
	LargeSeqUpdate bool   `jce:"14" json:"large_seq_update,omitempty"`

	Body0x769                []byte `jce:"15" json:"body_0x796,omitempty"`
	Status                   int32  `jce:"16" json:"status,omitempty"`
	ExtraStatus              int64  `jce:"17" json:"extra_status,omitempty"`
	ClientBatteryGetInterval int64  `jce:"18" json:"client_battery_get_interval,omitempty"`
	ClientAutoStatusInterval int64  `jce:"19" json:"client_auto_status_interval,omitempty"`
	CustomStatus             []byte `jce:"21" json:"custom_status,omitempty"`
}

type Register struct {
	Uin           int64    `jce:"1" json:"uin,omitempty"`
	PushIDs       []uint64 `jce:"2" json:"push_ids,omitempty"` // constant
	Status        uint32   `jce:"3" json:"status,omitempty"`
	KickPC        bool     `jce:"4" json:"kick_pc,omitempty"`
	KickWeak      bool     `jce:"5" json:"kick_weak,omitempty"` // constant false
	Timestamp     uint64   `jce:"6" json:"timestamp,omitempty"`
	LargeSeq      uint32   `jce:"7" json:"large_seq,omitempty"` // constant 0x00000000
	ExtraStatus   int64    `jce:"8" json:"extra_status,omitempty"`
	BatteryCap    int64    `jce:"9" json:"battery_cap,omitempty"`
	PowerConnect  bool     `jce:"10" json:"power_connect,omitempty"`
	BindUinNotify int64    `jce:"11" json:"bind_uin_notify,omitempty"`
}

func (m *Manager) Register(uin int64) {
}

func (m *Manager) pushRegister() (any, error) {
	buf, err := uni.Marshal(&uni.Data{
		Version:     0x0003,
		PacketType:  0x00,
		MessageType: 0x00000000,
		RequestID:   m.c.GetNextSeq(),
		ServantName: "PushService",
		FuncName:    "SvcReqRegister",
		Payload:     []byte{},
		Timeout:     0x00000000,
		Context:     map[string]string{},
		Status:      map[string]string{},
	}, map[string]any{
		"SvcReqRegister": nil,
	})
	if err != nil {
		return nil, err
	}
	return buf, nil
}
