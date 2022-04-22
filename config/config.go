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
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/elap5e/penguin/pkg/log"
)

type Config struct {
	Basedir  string     `json:"basedir" yaml:"basedir"`
	Servers  []*Server  `json:"servers" yaml:"servers"`
	Logging  *Logging   `json:"logging" yaml:"logging"`
	Storage  *Storage   `json:"storage" yaml:"storage"`
	Plugins  []*Plugin  `json:"plugins" yaml:"plugins"`
	Accounts []*Account `json:"accounts" yaml:"accounts"`
	Database *Database  `json:"database" yaml:"database"`
}

type Server struct {
	Type    string `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
}

type Logging struct {
	Basedir string `json:"basedir,omitempty" yaml:"basedir"`
	Level   string `json:"level,omitempty" yaml:"level"`
}

type Storage struct {
	Cache *StorageBackend `json:"cache" yaml:"cache"`
	Media *StorageBackend `json:"media" yaml:"media"`
}

type StorageBackend struct {
	Engine string `json:"engine" yaml:"engine"`
	Config any    `json:"config" yaml:"config"`
}

type Plugin struct {
	Plugin string       `json:"engine" yaml:"engine"`
	Config PluginConfig `json:"config" yaml:"config"`
}

type PluginConfig = any

type Account struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password,omitempty" yaml:"password"`
}

type Database struct {
	Engine string          `json:"engine" yaml:"engine"`
	Config *DatabaseConfig `json:"config" yaml:"config"`
}

type DatabaseConfig struct {
	Path     string `json:"path,omitempty" yaml:"path"`
	Host     string `json:"host,omitempty" yaml:"host"`
	Port     string `json:"port,omitempty" yaml:"port"`
	Database string `json:"database,omitempty" yaml:"database"`
	Username string `json:"username,omitempty" yaml:"username"`
	Password string `json:"password,omitempty" yaml:"password"`
}

const defaultConfig = `# Configuration for Penguin

basedir: .penguin
servers:
  - type: http
    address: localhost:6748
logging:
  level: debug
storage:
  cache:
    engine: fs
  media:
    engine: fs
plugins:
  - plugin: default
    config: {}
  - plugin: pprof
    config:
      address: localhost:6060
accounts:
  - username: 10000
    password: password
database:
  engine: default
`

func OpenFile(name string) *Config {
	data, err := ioutil.ReadFile(name)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path.Dir(name), 0750); err != nil && !os.IsExist(err) {
			log.Fatal("unable to create config directory: %s", err)
		}
		dirs := []string{
			".penguin/log",
			".penguin/cache/temps",
			".penguin/cache/service",
			".penguin/cache/session",
			".penguin/cache/tickets",
			".penguin/media",
			".penguin/media/blobs",
			".penguin/media/blobs/md5",
			".penguin/media/blobs/sha256",
			".penguin/media/metas",
			".penguin/media/audio",
			".penguin/media/document",
			".penguin/media/photo",
			".penguin/media/video",
			".penguin/media/voice",
		}
		oldDirs := []string{
			".penguin/cache",
			".penguin/cache/blobs",
			".penguin/cache/blobs/md5",
			".penguin/cache/blobs/sha256",
			".penguin/cache/metas",
			".penguin/cache/temps",
			".penguin/cache/audio",
			".penguin/cache/photo",
			".penguin/cache/video",
			".penguin/cache/voice",
			".penguin/service",
			".penguin/session",
			".penguin/tickets",
		}
		for _, dir := range append(dirs, oldDirs...) {
			if err = os.MkdirAll(dir, 0750); err != nil && !os.IsExist(err) {
				log.Fatal("unable to create directory: %s", err)
			}
		}
		log.Warn("config file does not exist, using default config: %s", name)
		data = []byte(defaultConfig)
		err = ioutil.WriteFile(name, data, 0640)
		if err != nil {
			log.Fatal("failed to write default config: %s", err)
		}
		log.Warn("please edit the config file and restart")
		os.Exit(0)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("failed to parse config: %s", err)
	}
	return setDefault(&cfg)
}

func setDefault(cfg *Config) *Config {
	if cfg.Basedir == "" {
		cfg.Basedir = ".penguin"
	}
	if cfg.Logging == nil {
		cfg.Logging = &Logging{
			Basedir: path.Join(cfg.Basedir, "log"),
			Level:   "debug",
		}
	} else if cfg.Logging.Basedir == "" {
		cfg.Logging.Basedir = path.Join(cfg.Basedir, "log")
	}
	return cfg
}
