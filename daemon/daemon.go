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

package deamon

import (
	"github.com/elap5e/penguin/daemon/account"
	"github.com/elap5e/penguin/daemon/auth"
	"github.com/elap5e/penguin/daemon/contact"
	"github.com/elap5e/penguin/daemon/message"
	"github.com/elap5e/penguin/pkg/net/msf/rpc"
)

type Daemon struct {
	c rpc.Client

	accm *account.Manager
	athm *auth.Manager
	cntm *contact.Manager
	msgm *message.Manager
}

func New() {

}