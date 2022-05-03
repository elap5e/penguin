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

type Source struct {
	App    *App
	SDK    *SDK
	Device *Device
}

func NewMobileQQAndroidSource(uin int64) *Source {
	return &Source{
		App:    MobileQQAndroidApp,
		SDK:    MobileQQAndroidSDK,
		Device: NewAndroidDevice(uin),
	}
}

func NewMiniHDQQAndroidSource(uin int64) *Source {
	return &Source{
		App:    MiniHDQQAndroidApp,
		SDK:    MiniHDQQAndroidSDK,
		Device: NewAndroidDevice(uin),
	}
}

func NewNextMobileQQAndroidSource(uin int64) *Source {
	return &Source{
		App:    NextMobileQQAndroidApp,
		SDK:    NextMobileQQAndroidSDK,
		Device: NewAndroidDevice(uin),
	}
}

func NewNextMiniHDQQAndroidSource(uin int64) *Source {
	return &Source{
		App:    NextMiniHDQQAndroidApp,
		SDK:    NextMiniHDQQAndroidSDK,
		Device: NewAndroidDevice(uin),
	}
}
