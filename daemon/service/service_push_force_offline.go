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
	"strconv"
	"time"

	"github.com/elap5e/penguin/pkg/encoding/uni"
	"github.com/elap5e/penguin/pkg/log"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type PushForceOfflineRequest struct {
	Uin        int64  `jce:"0" json:"uin"`
	Title      string `jce:"1" json:"title"`
	Message    string `jce:"2" json:"message"`
	SameDevice bool   `jce:"3" json:"same_device"`
}

func (m *Manager) handlePushForceOffline(reply *rpc.Reply) (*rpc.Args, error) {
	data, push := uni.Data{}, PushForceOfflineRequest{}
	if err := uni.Unmarshal(reply.Payload, &data, map[string]any{
		"req_PushForceOffline": &push,
	}); err != nil {
		return nil, err
	}
	return m.onPushForceOffline(push.Uin, &push)
}

func (m *Manager) onPushForceOffline(uin int64, push *PushForceOfflineRequest) (*rpc.Args, error) {
	log.Warn("service.onPushForceOffline uin:%d title:%q message:%q same_device:%t", push.Uin, push.Title, push.Message, push.SameDevice)
	duration := time.Second * 5
	timer, ok := m.timers[uin]
	if !ok {
		m.timers[uin] = time.NewTimer(duration)
		timer = m.timers[uin]
		go func() {
			for {
				<-timer.C
				if _, err := m.GetAuthManager().SignInCreateToken(strconv.FormatInt(uin, 10), uin); err != nil {
					log.Error("service.onPushForceOffline uin:%d SignInCreateToken error:%v", uin, err)
				}
				if _, err := m.RegisterAppRegister(uin); err != nil {
					log.Error("service.onPushForceOffline uin:%d RegisterAppRegister error:%v", uin, err)
				}
			}
		}()
	}
	timer.Reset(duration)
	log.Warn("service.onPushForceOffline uin:%d timer.Reset duration:%v", uin, duration)
	return nil, nil
}
