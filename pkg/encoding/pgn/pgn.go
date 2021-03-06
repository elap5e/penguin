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

package pgn

import (
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UTC().UnixMicro()))

func ParsePhotoType(ext string) uint32 {
	switch ext {
	case ".jpeg", ".jpg":
		return 1000
	case ".png":
		return 1001
	case ".webp":
		return 1002
	case ".sharpp":
		return 1004
	case ".bmp":
		return 1005
	case ".gif":
		return 2000
	case ".apng":
		return 2001
	}
	return 0
}
