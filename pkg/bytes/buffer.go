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

package bytes

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{Buffer: bytes.NewBuffer(buf)}
}

func (b *Buffer) ReadBool() (bool, error) {
	n, err := b.ReadByte()
	if err != nil {
		return false, err
	}
	return n != 0, nil
}

func (b *Buffer) ReadInt16() (int16, error) {
	n, err := b.ReadUint16()
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}

func (b *Buffer) ReadUint16() (uint16, error) {
	p := make([]byte, 2)
	if _, err := b.Read(p); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(p), nil
}

func (b *Buffer) ReadInt32() (int32, error) {
	n, err := b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

func (b *Buffer) ReadUint32() (uint32, error) {
	p := make([]byte, 4)
	if _, err := b.Read(p); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(p), nil
}

func (b *Buffer) ReadBytesL16V() ([]byte, error) {
	n, err := b.ReadUint16()
	if err != nil {
		return nil, err
	}
	p := make([]byte, n)
	if _, err := b.Read(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (b *Buffer) ReadBytesL16() ([]byte, error) {
	n, err := b.ReadUint16()
	if err != nil {
		return nil, err
	}
	p := make([]byte, n-2)
	if _, err := b.Read(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (b *Buffer) ReadBytesL32() ([]byte, error) {
	n, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	p := make([]byte, n-4)
	if _, err := b.Read(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (b *Buffer) ReadStringL16V() (string, error) {
	n, err := b.ReadUint16()
	if err != nil {
		return "", err
	}
	p := make([]byte, n)
	if _, err := b.Read(p); err != nil {
		return "", err
	}
	return string(p), nil
}

func (b *Buffer) ReadStringL16() (string, error) {
	n, err := b.ReadUint16()
	if err != nil {
		return "", err
	}
	p := make([]byte, n-2)
	if _, err := b.Read(p); err != nil {
		return "", err
	}
	return string(p), nil
}

func (b *Buffer) ReadStringL32() (string, error) {
	n, err := b.ReadUint32()
	if err != nil {
		return "", err
	}
	p := make([]byte, n-4)
	if _, err := b.Read(p); err != nil {
		return "", err
	}
	return string(p), nil
}

func (b *Buffer) WriteBool(v bool) error {
	if v {
		return b.WriteByte(1)
	}
	return b.WriteByte(0)
}

func (b *Buffer) WriteInt16(v int16) (int, error) {
	return b.WriteUint16(uint16(v))
}

func (b *Buffer) WriteUint16(v uint16) (int, error) {
	p := make([]byte, 2)
	binary.BigEndian.PutUint16(p, v)
	return b.Write(p)
}

func (b *Buffer) WriteInt32(v int32) (int, error) {
	return b.WriteUint32(uint32(v))
}

func (b *Buffer) WriteUint32(v uint32) (int, error) {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, v)
	return b.Write(p)
}

func (b *Buffer) WriteUint32At(v uint32, n int) (int, error) {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, v)
	copy(reflect.ValueOf(b.Buffer).Elem().FieldByName("buf").Bytes()[:4], p)
	return 4, nil
}

func (b *Buffer) WriteInt64(v int64) (int, error) {
	return b.WriteUint64(uint64(v))
}

func (b *Buffer) WriteUint64(v uint64) (int, error) {
	p := make([]byte, 8)
	binary.BigEndian.PutUint64(p, v)
	return b.Write(p)
}

func (b *Buffer) WriteBytesL16V(s []byte, l ...int16) (int, error) {
	n := len(s)
	if len(l) > 0 {
		if n > int(l[0]) {
			n = int(l[0])
		}
	}
	b.WriteUint16(uint16(n))
	n, _ = b.Write(s[:n])
	return n + 2, nil
}

func (b *Buffer) WriteBytesL16(s []byte) (int, error) {
	n := len(s)
	n, _ = b.WriteUint16(uint16(n + 2))
	n, _ = b.Write(s)
	return n + 2, nil
}

func (b *Buffer) WriteBytesL32(s []byte) (int, error) {
	n := len(s)
	n, _ = b.WriteUint32(uint32(n + 4))
	n, _ = b.Write(s)
	return n + 4, nil
}

func (b *Buffer) WriteStringL16V(s string, l ...int16) (int, error) {
	n := len(s)
	if len(l) > 0 {
		if n > int(l[0]) {
			n = int(l[0])
		}
	}
	b.WriteUint16(uint16(n))
	n, _ = b.WriteString(s)
	return n + 2, nil
}

func (b *Buffer) WriteStringL16(s string) (int, error) {
	n := len(s)
	n, _ = b.WriteUint16(uint16(n + 2))
	n, _ = b.WriteString(s)
	return n + 2, nil
}

func (b *Buffer) WriteStringL32(s string) (int, error) {
	n := len(s)
	n, _ = b.WriteUint32(uint32(n + 4))
	n, _ = b.WriteString(s)
	return n + 4, nil
}
