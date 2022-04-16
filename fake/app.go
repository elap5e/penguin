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

package fake

type App struct {
	FixID int32
	AppID int32

	PkgName  string
	VerName  string
	Revision string
	SigHash  [16]byte

	BuildAt int64
	SDKVer  string
	SSOVer  uint32

	ImageType  uint8
	MiscBitMap uint32

	CanCaptcha bool
}

var (
	AndroidMobileQQApp = &App{
		FixID:      537116314,
		AppID:      537116314,
		PkgName:    "com.tencent.mobileqq",
		VerName:    "8.8.85", // 8.8.85.7685
		Revision:   "8.8.85.6dd84073",
		SigHash:    [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d},
		BuildAt:    1645432578,
		SDKVer:     "6.0.0.2497",
		SSOVer:     18,
		ImageType:  1,
		CanCaptcha: true,
	}
	AndroidMobileQQAppNext = &App{
		FixID:      537117844,
		AppID:      537117844,
		PkgName:    "com.tencent.mobileqq",
		VerName:    "8.8.88", // 8.8.88.7770
		Revision:   "8.8.88.d47ef8e2",
		SigHash:    [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d},
		BuildAt:    1645432578,
		SDKVer:     "6.0.0.2497",
		SSOVer:     18,
		ImageType:  1,
		CanCaptcha: true,
	}
)
