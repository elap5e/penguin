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

type Codec interface {
	Close() error

	ReadResponseHeader(*Response) error
	ReadResponseBody(*Reply) error
	WriteRequest(*Request, *Args) error
}

const (
	EncryptTypeNotNeedEncrypt = 0x00
	EncryptTypeEncryptByD2Key = 0x01
	EncryptTypeEncryptByZeros = 0x02
)

const (
	FlagNoCompression   = 0x00000000
	FlagZlibCompression = 0x00000001
)

const (
	VersionDefault = 0x0000000a
	VersionSimple  = 0x0000000b
)
