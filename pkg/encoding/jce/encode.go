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

package jce

import (
	"reflect"
	"strconv"

	"github.com/elap5e/penguin/pkg/bytes"
)

func Marshal(v interface{}, opts ...bool) ([]byte, error) {
	simple := false
	if len(opts) != 0 && opts[0] {
		simple = true
	}

	e := encoder{}
	err := e.marshal(v, simple)
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)

	return buf, nil
}

type encoder struct {
	bytes.Buffer
}

func (e *encoder) marshal(v interface{}, simple bool) error {
	if !simple {
		e.reflectValue(reflect.ValueOf(v), 0x00)
	} else {
		e.reflectValue(reflect.ValueOf(v), 0xff)
	}
	return nil
}

func (e *encoder) reflectValue(v reflect.Value, t uint8) {
	typeEncoder(v.Type())(e, v, t)
}

func (e *encoder) EncodeHead(v uint8, t uint8) {
	if t < 15 {
		e.WriteByte(v | t<<4)
	} else {
		e.WriteByte(v | 240)
		e.WriteByte(t)
	}
}

type encoderFunc func(e *encoder, v reflect.Value, t uint8)

func typeEncoder(t reflect.Type) encoderFunc {
	switch t.Kind() {
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Ptr:
		return newPtrEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.String:
		return stringEncoder
	case reflect.Float64, reflect.Float32:
		return floatEncoder
	case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int16, reflect.Int8:
		return intEncoder
	case reflect.Uint64, reflect.Uint32, reflect.Uint, reflect.Uint16, reflect.Uint8:
		return uintEncoder
	case reflect.Bool:
		return boolEncoder
	default:
		return nil
	}
}

func interfaceEncoder(e *encoder, v reflect.Value, t uint8) {
	e.reflectValue(v.Elem(), t)
}

type ptrEncoder struct {
	elemEnc encoderFunc
}

func (pe ptrEncoder) encode(e *encoder, v reflect.Value, t uint8) {
	pe.elemEnc(e, v.Elem(), t)
}

func newPtrEncoder(t reflect.Type) encoderFunc {
	enc := ptrEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

func bytesEncoder(e *encoder, v reflect.Value, t uint8) {
	b := v.Bytes()
	e.EncodeHead(0x0d, t)
	e.EncodeHead(0x00, 0x00)
	uintEncoder(e, reflect.ValueOf(uint32(len(b))), 0x00)
	e.Write(b)
}

type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se sliceEncoder) encode(e *encoder, v reflect.Value, t uint8) {
	se.arrayEnc(e, v, t)
}

func newSliceEncoder(t reflect.Type) encoderFunc {
	if t.Elem().Kind() == reflect.Uint8 {
		return bytesEncoder
	}
	enc := sliceEncoder{newArrayEncoder(t)}
	return enc.encode
}

type field struct {
	name  string
	tag   uint8
	index []int
	typ   reflect.Type

	encoder encoderFunc
}

type structFields struct {
	list     []field
	tagIndex map[uint8]int
}

func typeFields(t reflect.Type) structFields {
	current := []field{}
	next := []field{{typ: t}}

	var count, nextCount map[reflect.Type]int

	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				tag := sf.Tag.Get("jce")
				if tag == "-" {
					continue
				}
				name, tag := parseTag(tag)
				if tag == "-" {
					continue
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if name != "" || ft.Kind() != reflect.Struct {
					if name == "" {
						name = sf.Name
					}
					t, _ := strconv.ParseUint(tag, 10, 8)
					field := field{
						name:  name,
						tag:   uint8(t),
						index: index,
						typ:   ft,
					}
					fields = append(fields, field)
					if count[f.typ] > 1 {
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	for i := range fields {
		f := &fields[i]
		f.encoder = typeEncoder(typeByIndex(t, f.index))
	}
	tagIndex := make(map[uint8]int, len(fields))
	for i, field := range fields {
		tagIndex[field.tag] = i
	}
	return structFields{fields, tagIndex}
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

type structEncoder struct {
	fields structFields
}

func (se structEncoder) encode(e *encoder, v reflect.Value, t uint8) {
	if t != 0xff {
		e.EncodeHead(0x0a, t)
		defer e.EncodeHead(0x0b, 0x00)
	}
	for i := range se.fields.list {
		f := &se.fields.list[i]
		f.encoder(e, v.Field(i), f.tag)
	}
}

func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{fields: typeFields(t)}
	return se.encode
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae arrayEncoder) encode(e *encoder, v reflect.Value, t uint8) {
	e.EncodeHead(0x09, t)
	l := v.Len()
	uintEncoder(e, reflect.ValueOf(uint32(l)), 0x00)
	for i := 0; i < l; i++ {
		ae.elemEnc(e, v.Index(i), 0x00)
	}
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	enc := arrayEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

type mapEncoder struct {
	elemEnc encoderFunc
}

func (me mapEncoder) encode(e *encoder, v reflect.Value, t uint8) {
	ks := v.MapKeys()
	e.EncodeHead(0x08, t)
	uintEncoder(e, reflect.ValueOf(uint32(len(ks))), 0x00)
	for _, k := range ks {
		b := v.MapIndex(k)
		stringEncoder(e, k, 0x00)
		me.elemEnc(e, b, 0x01)
	}
}

func newMapEncoder(t reflect.Type) encoderFunc {
	me := mapEncoder{typeEncoder(t.Elem())}
	return me.encode
}

func stringEncoder(e *encoder, v reflect.Value, t uint8) {
	b := v.String()
	l := len(b)
	if l > 255 {
		e.EncodeHead(0x07, t)
		uintEncoder(e, reflect.ValueOf(uint32(len(b))), 0x00)
		e.WriteString(b)
		return
	}
	e.EncodeHead(0x06, t)
	e.WriteByte(uint8(l))
	e.WriteString(b)
}

func floatEncoder(e *encoder, v reflect.Value, t uint8) {
	k := v.Kind()
	b := v.Float()
	switch k {
	case reflect.Float64:
		e.EncodeHead(0x05, t)
		e.WriteFloat64(b)
	case reflect.Float32:
		e.EncodeHead(0x04, t)
		e.WriteFloat32(float32(b))
	}
}

func intEncoder(e *encoder, v reflect.Value, t uint8) {
	k := v.Kind()
	b := v.Int()
	switch k {
	case reflect.Int64:
		if b>>32 != 0 {
			e.EncodeHead(0x03, t)
			e.WriteInt64(b)
			return
		}
		fallthrough
	case reflect.Int32, reflect.Int:
		if b>>16 != 0 {
			e.EncodeHead(0x02, t)
			e.WriteInt32(int32(b))
			return
		}
		fallthrough
	case reflect.Int16:
		if b>>8 != 0 {
			e.EncodeHead(0x01, t)
			e.WriteInt16(int16(b))
			return
		}
		fallthrough
	case reflect.Int8:
		if b != 0 {
			e.EncodeHead(0x00, t)
			e.WriteByte(byte(b))
			return
		}
		e.EncodeHead(0x0c, t)
	}
}

func uintEncoder(e *encoder, v reflect.Value, t uint8) {
	k := v.Kind()
	b := v.Uint()
	switch k {
	case reflect.Uint64:
		if b>>32 != 0 {
			e.EncodeHead(0x03, t)
			e.WriteUint64(b)
			return
		}
		fallthrough
	case reflect.Uint32, reflect.Uint:
		if b>>16 != 0 {
			e.EncodeHead(0x02, t)
			e.WriteUint32(uint32(b))
			return
		}
		fallthrough
	case reflect.Uint16:
		if b>>8 != 0 {
			e.EncodeHead(0x01, t)
			e.WriteUint16(uint16(b))
			return
		}
		fallthrough
	case reflect.Uint8:
		if b != 0 {
			e.EncodeHead(0x00, t)
			e.WriteByte(byte(b))
			return
		}
		e.EncodeHead(0x0c, t)
	}
}

func boolEncoder(e *encoder, v reflect.Value, t uint8) {
	if v.Bool() {
		e.EncodeHead(0x00, t)
		e.WriteByte(0x01)
		return
	}
	e.EncodeHead(0x0c, t)
}
