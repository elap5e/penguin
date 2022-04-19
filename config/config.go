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

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/elap5e/penguin/pkg/log"
)

type Config struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func OpenFile(name string) *Config {
	data, err := ioutil.ReadFile(name)
	if os.IsNotExist(err) {
		data, err = json.MarshalIndent(defaultConfig, "", "  ")
		if err == nil {
			err = ioutil.WriteFile(name, data, 0644)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return &cfg
}

var defaultConfig = &Config{
	Username: "10000",
}
