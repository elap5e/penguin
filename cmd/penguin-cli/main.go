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

package main

import "fmt"

func main() {
	kv := map[string]int64{
		"WLOGIN_A2":       64,
		"WLOGIN_A5":       2,
		"WLOGIN_AQSIG":    2097152,
		"WLOGIN_D2":       262144,
		"WLOGIN_DA2":      33554432,
		"WLOGIN_LHSIG":    4194304,
		"WLOGIN_LSKEY":    512,
		"WLOGIN_OPENKEY":  16384,
		"WLOGIN_PAYTOKEN": 8388608,
		"WLOGIN_PF":       16777216,
		"WLOGIN_PSKEY":    1048576,
		"WLOGIN_PT4Token": 134217728,
		"WLOGIN_QRPUSH":   67108864,
		"WLOGIN_RESERVED": 16,
		"WLOGIN_SID":      524288,
		"WLOGIN_SIG64":    8192,
		"WLOGIN_SKEY":     4096,
		"WLOGIN_ST":       128,
		"WLOGIN_STWEB":    32,
		"WLOGIN_TOKEN":    32768,
		"WLOGIN_VKEY":     131072,
	}
	for k, v := range kv {
		fmt.Printf("%s = 0b%032b\n", k, v)
	}
	fmt.Printf("%s = 0b%032b\n", "MiscBitmap", 0x08f7ff7c)
}
