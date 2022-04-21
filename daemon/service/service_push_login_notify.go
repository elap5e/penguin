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

	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type PushLoginNotifyRequest struct {
	AppID       int64      `jce:"0" json:"app_id"`
	Status      byte       `jce:"1" json:"status"`
	Tablet      byte       `jce:"2" json:"tablet"`
	Platform    int32      `jce:"3" json:"platform"`
	Title       string     `jce:"4" json:"title"`
	Message     string     `jce:"5" json:"message"`
	ProductType int32      `jce:"6" json:"product_type"`
	ClientType  int32      `jce:"7" json:"client_type"`
	Instances   []Instance `jce:"8" json:"instances"`
}

type Instance struct {
	AppID       int64 `jce:"0" json:"app_id"`
	Tablet      byte  `jce:"1" json:"tablet"`
	Platform    int32 `jce:"2" json:"platform"`
	ProductType int32 `jce:"3" json:"product_type"`
	ClientType  int32 `jce:"4" json:"client_type"`
}

func (m *Manager) handlePushLoginNotify(reply *rpc.Reply) (*rpc.Args, error) {
	data, push := uni.Data{}, PushLoginNotifyRequest{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"SvcReqMSFLoginNotify": &push,
	}); err != nil {
		return nil, err
	}
	p, _ := json.Marshal(push)
	log.Debug("service.handlePushLoginNotify: %s", p)
	return nil, nil
}
