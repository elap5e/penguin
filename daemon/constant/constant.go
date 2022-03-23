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

package constant

var (
	SMSAppID     = uint64(0x0000000000000009)
	AppID        = uint64(0x0000000000000010)
	SubAppID     = uint64(0x000000002a9e5427)
	OpenSDKID    = uint64(0x000000005f5e1604)
	SubAppIDList = []uint64{0x000000005f5e10e2}
)

var (
	SigMap            = uint32(0x00ff32f2)
	SubSigMap         = uint32(0x00010400)
	CodecAppIDMapByte = map[int]uint8{0: 2, 1: 0, 2: 1, 3: 3}
)

var (
	LocaleID = uint32(0x00000804)
	Domains  = []string{
		"game.qq.com",
		"mail.qq.com",
		"qzone.qq.com",
		"qun.qq.com",
		"openmobile.qq.com",
		"tenpay.com",
		"connect.qq.com",
		"qqweb.qq.com",
		"office.qq.com",
		"ti.qq.com",
		"mma.qq.com",
		"docs.qq.com",
		"vip.qq.com",
		"gamecenter.qq.com",
	}
)
