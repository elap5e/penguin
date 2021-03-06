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

	"github.com/elap5e/penguin/config"
	"github.com/elap5e/penguin/daemon"
	"github.com/elap5e/penguin/pkg/log"
)

func main() {
	config := config.OpenFile(".penguin/config.yaml")
	penguind := daemon.New(context.Background(), config)
	if err := penguind.Run(); err != nil {
		log.Error("penguin daemon exit with error: %s", err)
	}
}
