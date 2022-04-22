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
	"os"
	"path"

	"github.com/elap5e/penguin/config"
	"github.com/elap5e/penguin/daemon"
	"github.com/elap5e/penguin/pkg/log"
)

func main() {
	home, _ := os.UserHomeDir()
	_ = os.MkdirAll(path.Join(home, ".penguin"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "blobs"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "blobs", "md5"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "blobs", "sha256"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "metas"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "temps"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "audio"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "photo"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "video"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "cache", "voice"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "log"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "service"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "session"), 0755)
	_ = os.MkdirAll(path.Join(home, ".penguin", "tickets"), 0755)
	c := config.OpenFile(path.Join(home, ".penguin/config.yaml"))
	d := daemon.New(context.Background(), c)
	if err := d.Run(); err != nil {
		log.Error("penguin daemon exit with error: %s", err)
	}
}
