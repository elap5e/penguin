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

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func NewAndroidDevice(uin int64) *rpc.FakeDevice {
	r := rand.New(rand.NewSource(uin))
	buf := make([]byte, 20)
	_, err := r.Read(buf)
	if err != nil {
		log.Fatalf("failed to generate device config")
	}
	ipv4 := net.IPv4(192, 168, 0, buf[0])
	mac1 := fmt.Sprintf("00:50:%02x:%02x:00:%02x", buf[1], buf[2], buf[0])
	mac2 := fmt.Sprintf("00:50:%02x:%02x:00:%02x", buf[1], buf[2], buf[3])
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", buf[4:7], buf[8:9], buf[10:11], buf[12:13], buf[14:19])
	imei := fmt.Sprintf("86030802%07d", r.Int31n(10000000))
	imsi := fmt.Sprintf("460001%09d", r.Int31n(1000000000))
	osid := fmt.Sprintf("RKQ1.%07d.002", r.Int31n(10000000))
	return &rpc.FakeDevice{
		OS: rpc.FakeDeviceOS{
			Type:        "android",
			Version:     "11",
			BuildBrand:  []byte("Xiaomi"),
			BuildModel:  "Redmi K20",
			BuildID:     osid,
			SDKVersion:  30,
			NetworkType: 2,
		},
		APNName:       []byte("wifi"),
		SIMOPName:     []byte("CMCC"),
		Bootloader:    "unknown",
		ProcVersion:   "Linux version 2.6.18-92.el5 (brewbuilder@ls20-bc2-13.build.redhat.com)",
		Codename:      "davinci",
		Incremental:   "20.10.20",
		Fingerprint:   "Xiaomi/davinci/davinci:11/" + osid + "/20.10.20:user/release-keys",
		BootID:        uuid,
		Baseband:      "4.3.c5-00069-SM6150_GEN_PACK-1",
		InnerVersion:  "20.10.20",
		NetworkType:   1,
		NetIPFamily:   3,
		IPv4Address:   ipv4,
		IPv6Address:   ipv4,
		MACAddress:    mac1,
		BSSIDAddress:  mac2,
		SSIDAddress:   "unknown",
		IMEI:          imei,
		IMSI:          imsi,
		GUID:          md5.Sum(append([]byte(osid), mac1...)),
		GUIDFlag:      uint32((1 << 24 & 0xff000000) | (0 << 8 & 0xff00)),
		IsGUIDFileNil: false,
		IsGUIDGenSucc: true,
		IsGUIDChanged: false,
	}
}
