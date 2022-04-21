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

package oidb

import (
	"google.golang.org/protobuf/proto"

	"github.com/elap5e/penguin/pkg/encoding/oidb/pb"
)

type Data struct {
	Command uint32 `json:"command"`
	Service uint32 `json:"service"`
	Payload []byte `json:"payload"`
	Message string `json:"message"`
	Client  string `json:"client"`
	Result  uint32 `json:"result"`
}

func Marshal(v *Data) ([]byte, error) {
	return proto.Marshal(&pb.OidbSso_OIDBSSOPkg{
		Bodybuffer:    v.Payload,
		ClientVersion: v.Client,
		Command:       v.Command,
		ServiceType:   v.Service,
	})
}

func Unmarshal(b []byte, v *Data) error {
	var pkg pb.OidbSso_OIDBSSOPkg
	if err := proto.Unmarshal(b, &pkg); err != nil {
		return err
	}
	v.Command = pkg.GetCommand()
	v.Service = pkg.GetServiceType()
	v.Payload = append([]byte{}, pkg.GetBodybuffer()...)
	v.Message = pkg.GetErrorMsg()
	v.Client = pkg.GetClientVersion()
	v.Result = pkg.GetResult()
	return nil
}
