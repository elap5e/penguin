#!/usr/bin/env bash

# Copyright 2022 Elapse and contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -uex

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-darwin-amd64 cmd/penguin-cli/main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-darwin-arm64 cmd/penguin-cli/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-linux-386 cmd/penguin-cli/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-linux-amd64 cmd/penguin-cli/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-linux-arm64 cmd/penguin-cli/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-windows-386.exe cmd/penguin-cli/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-windows-amd64.exe cmd/penguin-cli/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o build/releases/penguin-cli-windows-arm64.exe cmd/penguin-cli/main.go

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-darwin-amd64 cmd/penguind/main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-darwin-arm64 cmd/penguind/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-linux-386 cmd/penguind/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-linux-amd64 cmd/penguind/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-linux-arm64 cmd/penguind/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-windows-386.exe cmd/penguind/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-windows-amd64.exe cmd/penguind/main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o build/releases/penguind-windows-arm64.exe cmd/penguind/main.go

upx -9 build/releases/penguin-cli-linux-{386,amd64,arm64} build/releases/penguin-cli-windows-{386,amd64}.exe
upx -9 build/releases/penguind-linux-{386,amd64,arm64} build/releases/penguind-windows-{386,amd64}.exe
