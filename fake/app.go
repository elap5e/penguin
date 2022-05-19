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

	Package  string
	Version  string
	Build    string
	Revision string
	SigHash  [16]byte
}

var (
	MobileQQAndroidApp = &App{
		FixID:    537119544,
		AppID:    537119544,
		Package:  "com.tencent.mobileqq",
		Version:  "8.8.90",
		Build:    "7975",
		Revision: "83e6c009",
		SigHash:  [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d},
	}
	MiniHDQQAndroidApp = &App{
		FixID:    537066967,
		AppID:    537066967,
		Package:  "com.tencent.minihd.qq",
		Version:  "5.9.4",
		Build:    "3666",
		Revision: "5ad348",
		SigHash:  [16]byte{0xaa, 0x39, 0x78, 0xf4, 0x1f, 0xd9, 0x6f, 0xf9, 0x91, 0x4a, 0x66, 0x9e, 0x18, 0x64, 0x74, 0xc7},
	}
)

var (
	NextMobileQQAndroidApp = MobileQQAndroidApp
	NextMiniHDQQAndroidApp = MiniHDQQAndroidApp
)
