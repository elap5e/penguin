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

package auth

import (
	"github.com/elap5e/penguin/daemon/constant"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
)

func (m *Manager) UnlockDevice(uin int64) (*Response, error) {
	fake, sess := m.c.GetFakeSource(uin), m.c.GetSession(uin)
	tlvs := make(map[uint16]tlv.Codec)
	tlvs[0x0008] = tlv.NewT8(0, constant.LocaleID, 0)
	tlvs[0x0104] = tlv.NewT104(sess.Auth)
	tlvs[0x0116] = tlv.NewT116(fake.App.MiscBitMap, constant.SubSigMap, constant.SubAppIDList)
	tlvs[0x0401] = tlv.NewT401(m.GetExtraData(uin).T401.Get())
	return m.requestSignIn(0, uin, 0x0014, tlvs)
}
