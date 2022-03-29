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
	"encoding/hex"
	"reflect"
	"testing"
)

func TestJCEEncoding(t *testing.T) {
	type testJCE struct {
		A bool              `jce:"1"`
		B uint8             `jce:"2"`
		C uint16            `jce:"3"`
		D uint32            `jce:"4"`
		E uint64            `jce:"5"`
		F float32           `jce:"6"`
		G float64           `jce:"7"`
		H string            `jce:"8"`
		I map[string][]byte `jce:"9"`
		J map[string]string `jce:"10"`
		K []uint64          `jce:"11"`
		L []byte            `jce:"12"`
	}
	lstr := "ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP"
	type args struct {
		v    any
		opts []bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "0x00",
			args: args{
				v: testJCE{
					A: false,
					B: 0x00,
					C: 0x0000,
					D: 0x00000000,
					E: 0x0000000000000000,
					F: 0x00000000,
					G: 0x0000000000000000,
					H: "",
					I: map[string][]byte{"": []byte("")},
					J: map[string]string{"": ""},
					K: []uint64{0x0000000000000000},
					L: []byte(""),
				},
				opts: []bool{false},
			},
			wantErr: false,
		},
		{
			name: "0x01",
			args: args{
				v: testJCE{
					A: true,
					B: 0x20,
					C: 0x3000,
					D: 0x40000000,
					E: 0x5000000000000000,
					F: 0x60000000,
					G: 0x7000000000000000,
					H: "H",
					I: map[string][]byte{"I": []byte("I")},
					J: map[string]string{"J": "J"},
					K: []uint64{0xb000000000000000},
					L: []byte("L"),
				},
				opts: []bool{false},
			},
			wantErr: false,
		},
		{
			name: "0x02",
			args: args{
				v: testJCE{
					A: true,
					B: 0xff,
					C: 0xffff,
					D: 0xffffffff,
					E: 0xffffffff00000000,
					F: 0xffffffff,
					G: 0xffffffff00000000,
					H: lstr,
					I: map[string][]byte{lstr: []byte("I")},
					J: map[string]string{lstr: lstr},
					K: []uint64{0xffffffff00000000},
					L: []byte(lstr),
				},
				opts: []bool{false},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data1, data2 []byte
			var err error
			if data1, err = Marshal(tt.args.v, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			v := new(testJCE)
			if err = Unmarshal(data1, v, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if data2, err = Marshal(v, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(data2, data1) {
				t.Errorf("Marshal() = \n%s, want\n%s", hex.Dump(data2), hex.Dump(data1))
			}
		})
	}
}
