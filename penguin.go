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

//go:generate protoc --go_out=$GOPATH/src daemon/auth/pb/gateway_verify.proto
//go:generate protoc --go_out=$GOPATH/src daemon/auth/pb/third_part_login.proto

//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/body.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/head.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/common.proto

//go:generate protoc --go_out=$GOPATH/src daemon/service/pb/domain_ip.proto
//go:generate protoc --go_out=$GOPATH/src daemon/service/pb/oidb_0x769.proto
//go:generate protoc --go_out=$GOPATH/src daemon/service/pb/online_push.proto

//go:generate protoc --go_out=$GOPATH/src pkg/encoding/tlv/pb/device_report.proto

//go:generate go run cmd/proto-message-gen/main.go

package penguin
