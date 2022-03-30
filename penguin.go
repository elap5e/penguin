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

//go:generate go run cmd/message-face-gen/main.go
//go:generate go run cmd/proto-message-gen/main.go

//go:generate protoc --go_out=$GOPATH/src pkg/encoding/tlv/pb/device_report.proto

//go:generate protoc --go_out=$GOPATH/src daemon/auth/pb/gateway_verify.proto
//go:generate protoc --go_out=$GOPATH/src daemon/auth/pb/third_part_login.proto

//go:generate protoc --go_out=$GOPATH/src daemon/contact/pb/mutual_mark.proto
//go:generate protoc --go_out=$GOPATH/src daemon/contact/pb/oidb_0xd50.proto
//go:generate protoc --go_out=$GOPATH/src daemon/contact/pb/oidb_0xd6b.proto

//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/body.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/head.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/common.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/common_element.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/control.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/receipt.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/sub_type_0x1a.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/sub_type_0xc1.proto
//go:generate protoc --go_out=$GOPATH/src daemon/message/pb/service.proto

//go:generate protoc --go_out=$GOPATH/src daemon/service/pb/domain_ip.proto
//go:generate protoc --go_out=$GOPATH/src daemon/service/pb/oidb_0x769.proto
//go:generate protoc --go_out=$GOPATH/src daemon/service/pb/online_push.proto

//go:generate protoc --go_out=$GOPATH/src daemon/channel/pb/channel_base.proto
//go:generate protoc --go_out=$GOPATH/src daemon/channel/pb/channel_common.proto
//go:generate protoc --go_out=$GOPATH/src daemon/channel/pb/message_common.proto
//go:generate protoc --go_out=$GOPATH/src daemon/channel/pb/message_push.proto
//go:generate protoc --go_out=$GOPATH/src daemon/channel/pb/sync_logic.proto

package penguin
