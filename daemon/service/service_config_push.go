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
	"encoding/json"
	"io/ioutil"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/daemon/service/pb"
	"github.com/elap5e/penguin/pkg/encoding/jce"
	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

type ConfigPushType uint32

var (
	ConfigPushTypeSSOServerConfig   ConfigPushType = 1
	ConfigPushTypeFileStorageServer ConfigPushType = 2
	ConfigPushTypeClientLogConfig   ConfigPushType = 3
	ConfigPushTypeProxyIPChannel    ConfigPushType = 4
)

type ConfigPushRequest struct {
	Type    ConfigPushType `jce:"1" json:"type"`
	Seq     int64          `jce:"3" json:"seq"`
	Payload []byte         `jce:"2" json:"payload"`
}

type ConfigPushResponse struct {
	Type    ConfigPushType `jce:"1" json:"type"`
	Seq     int64          `jce:"2" json:"seq"`
	Payload []byte         `jce:"3" json:"payload"`
}

type SSOServerConfig struct {
	TGTGList     []*SSOServerInfo `jce:"1" json:"2g3g_list,omitempty"`
	WiFiList     []*SSOServerInfo `jce:"3" json:"wifi_list,omitempty"`
	Reconnect    int32            `jce:"4" json:"reconnect,omitempty"`
	TestSpeed    bool             `jce:"5" json:"test_speed,omitempty"`
	UseNewList   bool             `jce:"6" json:"use_new_list,omitempty"`
	MultiConnect int32            `jce:"7" json:"multi_connect,omitempty"`
	HTTP2G3GList []*SSOServerInfo `jce:"8" json:"http_2g3g_list,omitempty"`
	HTTPWiFiList []*SSOServerInfo `jce:"9" json:"http_wifi_list,omitempty"`
} // SsoServerList

type SSOServerInfo struct {
	IP           string `jce:"1" json:"ip,omitempty"`
	Port         int32  `jce:"2" json:"port,omitempty"`
	LinkType     bool   `jce:"3" json:"link_type,omitempty"`
	Proxy        bool   `jce:"4" json:"proxy,omitempty"`
	ProtocolType bool   `jce:"5" json:"protocol_type,omitempty"`
	Timeout      int32  `jce:"6" json:"timeout,omitempty"`
	Location     string `jce:"8" json:"location,omitempty"`
} // SsoServerListInfo

type FileStorageServer struct {
	UpLoadList               []*FileStorageServerInfo `jce:"0" json:"upload_list,omitempty"`
	PictureDownloadList      []*FileStorageServerInfo `jce:"1" json:"picture_download_list,omitempty"`
	GroupPictureDownloadList []*FileStorageServerInfo `jce:"2" json:"group_picture_download_list,omitempty"`
	QZoneProxyServerList     []*FileStorageServerInfo `jce:"3" json:"qzone_proxy_server_list,omitempty"`
	URLEncodeServerList      []*FileStorageServerInfo `jce:"4" json:"url_encode_server_list,omitempty"`
	BigDataIPChannel         *BigDataIPChannel        `jce:"5" json:"big_data_ip_channel,omitempty"`
	VIPEmotionList           []*FileStorageServerInfo `jce:"6" json:"vip_emotion_list,omitempty"`
	C2CPictureDownloadList   []*FileStorageServerInfo `jce:"7" json:"c2c_picture_download_list,omitempty"`
	FormatIPInfo             *FormatIPInfo            `jce:"8" json:"format_ip_info,omitempty"`
	DomainIPChannel          *DomainIPChannel         `jce:"9" json:"domain_ip_channel,omitempty"`
	PTTList                  []byte                   `jce:"10" json:"ptt_list,omitempty"`
} // FileStoragePushFSSvcList

type FileStorageServerInfo struct {
	IP   string `jce:"1" json:"ip,omitempty"`
	Port int32  `jce:"2" json:"port,omitempty"`
} // FileStorageServerListInfo

type FormatIPInfo struct {
	IP       string `jce:"0" json:"ip,omitempty"`
	Operator int64  `jce:"1" json:"operator,omitempty"`
}

type BigDataIPChannel struct {
	BigDataIPList []*BigDataIP `jce:"0" json:"big_data_ip_list,omitempty"`
	Sig           []byte       `jce:"1" json:"sig,omitempty"`
	Key           []byte       `jce:"2" json:"key,omitempty"`
	Uin           int64        `jce:"3" json:"uin,omitempty"`
	Flag          uint32       `jce:"4" json:"flag,omitempty"`
	Payload       []byte       `jce:"5" json:"payload,omitempty"`
}

type BigDataIP struct {
	Type       int64            `jce:"0" json:"type,omitempty"`
	IPList     []*BigDataIPInfo `jce:"1" json:"ip_list,omitempty"`
	ConfigList []*NetSegConfig  `jce:"2" json:"config_list,omitempty"`
	Size       int64            `jce:"3" json:"size,omitempty"`
}

type BigDataIPInfo struct {
	Type int64  `jce:"0" json:"type,omitempty"`
	IP   string `jce:"1" json:"ip,omitempty"`
	Port int64  `jce:"2" json:"port,omitempty"`
}

type NetSegConfig struct {
	NetType           int64 `jce:"0" json:"net_type,omitempty"`
	SegSize           int64 `jce:"1" json:"seg_size,omitempty"`
	SegNumber         int64 `jce:"2" json:"seg_number,omitempty"`
	CurrentConnNumber int64 `jce:"3" json:"current_conn_number,omitempty"`
}

type DomainIPChannel struct {
	DomainIPList []*DomainIP `jce:"0" json:"domain_ip_list,omitempty"`
}

type DomainIP struct {
	Type   uint32          `jce:"0" json:"type,omitempty"`
	IPList []*DomainIPInfo `jce:"1" json:"ip_list,omitempty"`
}

type DomainIPInfo struct {
	IP   uint32 `jce:"1" json:"ip,omitempty"`
	Port uint32 `jce:"2" json:"port,omitempty"`
}

type ClientLogConfig struct {
	Type       uint32     `jce:"1" json:"type,omitempty"`
	TimeStart  *Timestamp `jce:"2" json:"time_start,omitempty"`
	TimeFinish *Timestamp `jce:"3" json:"time_finish,omitempty"`
	LogLevel   uint8      `jce:"4" json:"log_level,omitempty"`
	Cookie     uint32     `jce:"5" json:"cookie,omitempty"`
	Seq        int64      `jce:"6" json:"seq,omitempty"`
}

type Timestamp struct {
	Year  uint32 `jce:"1" json:"year,omitempty"`
	Month uint8  `jce:"2" json:"month,omitempty"`
	Day   uint8  `jce:"3" json:"day,omitempty"`
	Hour  uint8  `jce:"4" json:"hour,omitempty"`
}

type ProxyIPChannel struct {
	ProxyIPList []*ProxyIP `jce:"0" json:"proxy_ip_list,omitempty"`
	Reconnect   uint32     `jce:"1" json:"reconnect,omitempty"`
}

type ProxyIP struct {
	Type   int64          `jce:"0" json:"type,omitempty"`
	IPlist []*ProxyIPInfo `jce:"1" json:"ip_list,omitempty"`
}

type ProxyIPInfo struct {
	Type uint32 `jce:"0" json:"type,omitempty"`
	IP   uint32 `jce:"1" json:"ip,omitempty"`
	Port uint32 `jce:"2" json:"port,omitempty"`
}

func (m *Manager) handleConfigPushDomain(reply *rpc.Reply) (*rpc.Args, error) {
	push := pb.DomainIp_NameRspBody{}
	if err := proto.Unmarshal(reply.Payload, &push); err != nil {
		return nil, err
	}
	file := ".penguin/service/config_push_domain." + strconv.FormatInt(reply.Uin, 10) + ".json"
	data, err := json.MarshalIndent(push.GetSubCmdName_Rsp(), "", "  ")
	if err == nil {
		err = ioutil.WriteFile(file, data, 0644)
	}
	return nil, nil
}

func (m *Manager) handleConfigPushRequest(reply *rpc.Reply) (*rpc.Args, error) {
	data, req := uni.Data{}, ConfigPushRequest{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"PushReq": &req,
	}); err != nil {
		return nil, err
	}
	switch req.Type {
	case ConfigPushTypeSSOServerConfig:
		cfg := SSOServerConfig{}
		if err := jce.Unmarshal(req.Payload, &cfg, true); err != nil {
			return nil, err
		}
		file := ".penguin/service/config_push_sso_server_config." + strconv.FormatInt(reply.Uin, 10) + ".json"
		data, err := json.MarshalIndent(&cfg, "", "  ")
		if err == nil {
			err = ioutil.WriteFile(file, data, 0644)
		}
		// TODO: update sso server config
	case ConfigPushTypeFileStorageServer:
		cfg := FileStorageServer{}
		if err := jce.Unmarshal(req.Payload, &cfg, true); err != nil {
			return nil, err
		}
		file := ".penguin/service/config_push_file_storage_server." + strconv.FormatInt(reply.Uin, 10) + ".json"
		data, err := json.MarshalIndent(&cfg, "", "  ")
		if err == nil {
			err = ioutil.WriteFile(file, data, 0644)
		}
		// TODO: update file storage server
	case ConfigPushTypeClientLogConfig:
		fallthrough
	case ConfigPushTypeProxyIPChannel:
		fallthrough
	default:
	}
	resp := &ConfigPushResponse{
		Type:    req.Type,
		Seq:     req.Seq,
		Payload: req.Payload,
	}
	p, err := uni.Marshal(&uni.Data{
		Version:     3,
		ServantName: data.ServantName,
		FuncName:    "PushResp",
	}, map[string]any{
		"PushResp": resp,
	})
	if err != nil {
		return nil, err
	}
	return &rpc.Args{
		Version:       rpc.VersionSimple,
		Uin:           reply.Uin,
		Seq:           reply.Seq,
		ServiceMethod: service.MethodServiceConfigPushResponse,
		Payload:       p,
	}, nil
}
