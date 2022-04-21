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

package msf

import (
	"crypto/rand"

	"github.com/elap5e/penguin/pkg/crypto/ecdh"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

func newSession() *rpc.Session {
	var session rpc.Session
	session = rpc.Session{}
	session.Cookie = make([]byte, 4)
	rand.Read(session.Cookie)
	session.RandomKey = [16]byte{}
	rand.Read(session.RandomKey[:])
	session.RandomPass = [16]byte{}
	rand.Read(session.RandomPass[:])
	strs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range session.RandomPass {
		session.RandomPass[i] = strs[session.RandomPass[i]%52]
	}
	session.PrivateKey, _ = ecdh.GenerateKey()
	session.KeyVersion = ecdh.ServerKeyVersion
	session.SharedSecret = session.PrivateKey.SharedSecret(ecdh.ServerPublicKey)
	return &session
}

func getSession(uin int64) *rpc.Session {
	return newSession()
}
