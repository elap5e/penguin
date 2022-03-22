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

package tea

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

type Cipher struct {
	key []uint32
	tmp [][]byte
}

func NewCipher(v [16]byte) *Cipher {
	key := make([]uint32, 4)
	key[0] = binary.BigEndian.Uint32(v[0:])
	key[1] = binary.BigEndian.Uint32(v[4:])
	key[2] = binary.BigEndian.Uint32(v[8:])
	key[3] = binary.BigEndian.Uint32(v[12:])
	var tmp = make([][]byte, 4)
	tmp[0] = make([]byte, 8)
	tmp[1] = make([]byte, 8)
	tmp[2] = make([]byte, 8)
	tmp[3] = make([]byte, 8)
	return &Cipher{key: key, tmp: tmp}
}

func (c *Cipher) encrypt(src []byte, off, n int) (dst []byte) {
	fill := 10 - (n-off+1)%8
	dst = make([]byte, fill+n-off+7)
	_, _ = rand.Read(dst[0:fill])
	dst[0] = (dst[0] & 0xF8) | byte(fill-3)
	copy(dst[fill:], src[off:off+n])
	for i := 0; i < len(dst); i += 8 {
		c.encrypt64bit(dst[i:i+8], dst[i:i+8])
	}
	return dst
}

func (c *Cipher) encrypt64bit(dst, src []byte) {
	xorBytes(dst, src, c.tmp[0])
	copy(c.tmp[0], dst)
	c.encryptBlock(dst, dst)
	xorBytes(dst, dst, c.tmp[1])
	copy(c.tmp[1], c.tmp[0])
	copy(c.tmp[0], dst)
}

func (c *Cipher) encryptBlock(dst, src []byte) {
	_ = src[7] // early bounds check
	s0 := binary.BigEndian.Uint32(src[0:4])
	s1 := binary.BigEndian.Uint32(src[4:8])
	var sum uint32 = 0x00000000
	for i := 0; i < 0x10; i++ {
		sum += 0x9e3779B9 // -1640531527
		s0 += ((s1 << 4) + c.key[0]) ^ (s1 + sum) ^ ((s1 >> 5) + c.key[1])
		s1 += ((s0 << 4) + c.key[2]) ^ (s0 + sum) ^ ((s0 >> 5) + c.key[3])
	}
	binary.BigEndian.PutUint32(dst[0:4], s0)
	binary.BigEndian.PutUint32(dst[4:8], s1)
}

func (c *Cipher) decrypt(src []byte, off, n int) (dst []byte, err error) {
	if (n-off)%8 != 0 {
		return nil, fmt.Errorf("invalid length")
	}
	dst = make([]byte, n-off)
	for i := 0; i < len(dst); i += 8 {
		c.decrypt64bit(dst[i:i+8], src[off+i:off+i+8])
	}
	for _, b := range dst[n-off-7:] {
		if b != 0x00 {
			return nil, fmt.Errorf("fail to decrypt")
		}
	}
	return dst[dst[0]&0x07+3 : n-off-7], nil
}

func (c *Cipher) decrypt64bit(dst, src []byte) {
	// TODO: fix src&dst memory overlap
	xorBytes(dst, src, c.tmp[2])
	c.decryptBlock(dst, dst)
	copy(c.tmp[2], dst)
	xorBytes(dst, dst, c.tmp[3])
	copy(c.tmp[3], src)
}

func (c *Cipher) decryptBlock(dst, src []byte) {
	_ = src[7] // early bounds check
	s0 := binary.BigEndian.Uint32(src[0:4])
	s1 := binary.BigEndian.Uint32(src[4:8])
	var sum uint32 = 0xE3779B90 // -478700656
	for i := 0; i < 0x10; i++ {
		s1 -= ((s0 << 4) + c.key[2]) ^ (s0 + sum) ^ ((s0 >> 5) + c.key[3])
		s0 -= ((s1 << 4) + c.key[0]) ^ (s1 + sum) ^ ((s1 >> 5) + c.key[1])
		sum -= 0x9e3779B9 // -1640531527
	}
	binary.BigEndian.PutUint32(dst[0:4], s0)
	binary.BigEndian.PutUint32(dst[4:8], s1)
}

func (c *Cipher) Encrypt(src []byte) []byte {
	return c.encrypt(src, 0, len(src))
}

func (c *Cipher) Decrypt(src []byte) ([]byte, error) {
	return c.decrypt(src, 0, len(src))
}

func (c *Cipher) SetKey(v [16]byte) {
	c.key[0] = binary.BigEndian.Uint32(v[0:])
	c.key[1] = binary.BigEndian.Uint32(v[4:])
	c.key[2] = binary.BigEndian.Uint32(v[8:])
	c.key[3] = binary.BigEndian.Uint32(v[12:])
}
