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

package auth

import (
	"crypto/sha256"
	"log"
	"math/big"
	"reflect"
	"time"

	"github.com/elap5e/penguin/pkg/bytes"
)

func calcPow(data []byte) []byte {
	tmp := bytes.NewBuffer(data)
	a, _ := tmp.ReadByte()
	typ, _ := tmp.ReadByte()
	c, _ := tmp.ReadByte()
	ok, _ := tmp.ReadBool()
	e, _ := tmp.ReadUint16()
	f := make([]byte, 2)
	_, _ = tmp.Read(f)
	src, _ := tmp.ReadBytesL16V()
	tgt, _ := tmp.ReadBytesL16V()
	cpy, _ := tmp.ReadBytesL16V()
	dst, elp, cnt := []byte{}, uint32(0), uint32(0)

	if typ == 2 {
		start := time.Now()
		tmp := new(big.Int).SetBytes(src)
		hash := sha256.Sum256(tmp.Bytes())
		for cnt = 0; !reflect.DeepEqual(hash[:], tgt); cnt++ {
			tmp = tmp.Add(tmp, big.NewInt(1))
			hash = sha256.Sum256(tmp.Bytes())
		}
		ok = true
		dst = tmp.Bytes()
		elp = uint32(time.Now().Sub(start).Milliseconds())
	} else {
		log.Fatalln("not support")
	}

	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(a)
	buf.WriteByte(typ)
	buf.WriteByte(c)
	buf.WriteBool(ok)
	buf.WriteUint16(e)
	buf.Write(f)
	buf.WriteBytesL16V(src)
	buf.WriteBytesL16V(tgt)
	buf.WriteBytesL16V(cpy)
	if ok {
		buf.WriteBytesL16V(dst)
		buf.WriteUint32(elp)
		buf.WriteUint32(cnt)
	}
	return buf.Bytes()
}
