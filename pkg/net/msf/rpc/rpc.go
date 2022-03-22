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

package rpc

type Call struct {
	ServiceMethod string
	Seq           uint64
	Args          *Args
	Reply         *Reply
	Error         error
	Done          chan *Call
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
	}
}

type Args struct {
	Seq    int32
	AppID  int32
	FixID  int32
	Buffer []byte
}

type Reply struct {
	Seq int32
}

type Request struct {
	ServiceMethod string
	Seq           uint64
}

type Response struct {
	ServiceMethod string
	Seq           uint64
}
