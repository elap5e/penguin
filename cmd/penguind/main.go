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

import (
	"context"
	"log"
	"math/rand"

	"github.com/elap5e/penguin/pkg/net/msf"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
	"github.com/elap5e/penguin/pkg/net/msf/service"
)

func main() {
	c := msf.NewClient(context.Background())
	log.Println(c.Call(service.MethodHeartbeatAlive, &rpc.Args{
		Seq:     rand.Int31n(100000),
		FixID:   537044845,
		AppID:   537044845,
		Payload: []byte{0x00, 0x00, 0x00, 0x04},
	}, nil))
}
