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

package uni

type Data struct {
	Version     int16             `jce:",1" json:"version,omitempty"`
	PacketType  uint8             `jce:",2" json:"packet_type,omitempty"`
	MessageType uint32            `jce:",3" json:"message_type,omitempty"`
	RequestID   int32             `jce:",4" json:"request_id,omitempty"`
	ServantName string            `jce:",5" json:"servant_name,omitempty"`
	FuncName    string            `jce:",6" json:"func_name,omitempty"`
	Payload     []byte            `jce:",7" json:"payload,omitempty"`
	Timeout     int32             `jce:",8" json:"timeout,omitempty"`
	Context     map[string]string `jce:",9" json:"context,omitempty"`
	Status      map[string]string `jce:",10" json:"status,omitempty"`
}

type Payload map[string]map[string][]byte

type PayloadV3 map[string][]byte
