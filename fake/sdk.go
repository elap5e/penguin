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

type SDK struct {
	BuildAt    int64
	Version    string
	SSOVersion uint32
	ImageType  uint8
	CanCaptcha bool
}

var (
	MobileQQAndroidSDK = &SDK{
		BuildAt:    1652271523,
		Version:    "6.0.0.2508",
		SSOVersion: 18,
		ImageType:  1,
		CanCaptcha: true,
	}
	MiniHDQQAndroidSDK = &SDK{
		BuildAt:    1630653275,
		Version:    "6.0.0.2484",
		SSOVersion: 18,
		ImageType:  1,
		CanCaptcha: true,
	}
)

var (
	NextMobileQQAndroidSDK = MobileQQAndroidSDK
	NextMiniHDQQAndroidSDK = MiniHDQQAndroidSDK
)
