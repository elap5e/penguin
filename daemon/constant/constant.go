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
	AppDomains = []string{
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
	SubAppIDList = []uint64{0x000000005f5e10e2}
	MapAppIDByte = map[int]uint8{0: 2, 1: 0, 2: 1, 3: 3}
)

const (
	LocaleID  = uint32(0x00000804)
	SMSAppID  = uint64(0x0000000000000009)
	DstAppID  = uint64(0x0000000000000010)
	OpenSDKID = uint64(0x000000005f5e1604)
)

// WtloginHelper
const (
	BitMapA5       = uint32(0b00000000000000000000000000000010)
	BitMapRESERVED = uint32(0b00000000000000000000000000010000)
	BitMapSTWEB    = uint32(0b00000000000000000000000000100000)
	BitMapA2       = uint32(0b00000000000000000000000001000000)
	BitMapST       = uint32(0b00000000000000000000000010000000)
	BitMapLSKey    = uint32(0b00000000000000000000001000000000)
	BitMapSKey     = uint32(0b00000000000000000001000000000000)
	BitMapSig64    = uint32(0b00000000000000000010000000000000)
	BitMapOpenKey  = uint32(0b00000000000000000100000000000000)
	BitMapToken    = uint32(0b00000000000000001000000000000000)
	BitMapVKey     = uint32(0b00000000000000100000000000000000)
	BitMapD2       = uint32(0b00000000000001000000000000000000)
	BitMapSID      = uint32(0b00000000000010000000000000000000)
	BitMapPSKey    = uint32(0b00000000000100000000000000000000)
	BitMapAQSig    = uint32(0b00000000001000000000000000000000)
	BitMapLHSig    = uint32(0b00000000010000000000000000000000)
	BitMapPAYToken = uint32(0b00000000100000000000000000000000)
	BitMapPF       = uint32(0b00000001000000000000000000000000)
	BitMapDA2      = uint32(0b00000010000000000000000000000000)
	BitMapQRPush   = uint32(0b00000100000000000000000000000000)
	BitMapPT4Token = uint32(0b00001000000000000000000000000000)
	MainSigMap     = uint32(0b00000000111111110011001011110010)
	SubSigMap      = uint32(0b00000000000000010000010000000000)
	OpenAppID      = uint64(0x000000002a9e5427)
)
