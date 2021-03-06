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
	// 0b00000000111111110011001011110010 16724722
	MainSigMap = SigType(0) |
		SigTypeA5 |
		SigTypeReserved |
		SigTypeSTWeb |
		SigTypeA2 |
		SigTypeST |
		SigTypeLSKey |
		SigTypeSKey |
		SigTypeSig64 |
		SigType(1<<16) |
		SigTypeVKey |
		SigTypeD2 |
		SigTypeSID |
		SigTypePSKey |
		SigTypeAQSig |
		SigTypeLHSig |
		SigTypePayToken
	// 0b00000000000000010000010000000000 66560
	SubSigMap = SigType(0) | SigType(1<<10) | SigType(1<<16)

	MiscBitMap = uint32(0b00001000111101111111111101111100) // 150470524
	OpenAppID  = uint64(0x000000002a9e5427)                 // 715019303
)

// WtloginHelper.SigType
type SigType = uint32

const (
	_ SigType = 1 << iota
	SigTypeA5
	_
	_
	SigTypeReserved
	SigTypeSTWeb
	SigTypeA2
	SigTypeST
	_
	SigTypeLSKey
	_
	_
	SigTypeSKey
	SigTypeSig64
	SigTypeOpenKey
	SigTypeToken
	_
	SigTypeVKey
	SigTypeD2
	SigTypeSID
	SigTypePSKey
	SigTypeAQSig
	SigTypeLHSig
	SigTypePayToken
	SigTypePF
	SigTypeDA2
	SigTypeQRPush
	SigTypePT4Token
	_
	_
	_
	_
)
