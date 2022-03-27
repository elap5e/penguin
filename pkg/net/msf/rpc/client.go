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

import (
	"encoding/hex"
	"net"
	"strconv"
	"time"

	"github.com/elap5e/penguin/pkg/crypto/ecdh"
	"github.com/elap5e/penguin/pkg/encoding/tlv"
)

type Client interface {
	Close() error
	Go(serviceMethod string, args *Args, reply *Reply, done chan *Call) *Call
	Call(serviceMethod string, args *Args, reply *Reply) error
	Handle(serviceMethod string, reply *Reply) (*Args, error)
	Register(serviceMethod string, handler Handler) error

	GetNextSeq() int32
	GetFakeSource(uin int64) *FakeSource
	GetServerTime() int64
	GetSession(uin int64) *Session
	GetTickets(uin int64) *Tickets
	SetSession(uin int64, tlvs map[uint16]tlv.Codec)
	SetSessionAuth(uin int64, auth []byte)
	SetSessionCookie(uin int64, cookie []byte)
	SetTickets(uin int64, tlvs map[uint16]tlv.Codec)
}

type Handler func(reply *Reply) (*Args, error)

type Session struct {
	Auth   []byte `json:"auth,omitempty"`
	Cookie []byte `json:"cookie,omitempty"`

	RandomKey  Key16Bytes `json:"random_key,omitempty"`
	RandomPass Key16Bytes `json:"random_pass,omitempty"`

	PrivateKey   *ecdh.PrivateKey `json:"-"`
	KeyVersion   int16            `json:"key_version,omitempty"`
	SharedSecret Key16Bytes       `json:"shared_secret,omitempty"`
}

type Tickets struct {
	A1       *Ticket           `json:"a1,omitempty"`
	A2       *Ticket           `json:"a2,omitempty"`
	A5       *Ticket           `json:"a5,omitempty"`
	A8       *Ticket           `json:"a8,omitempty"`
	D2       *Ticket           `json:"d2,omitempty"`
	LSKey    *Ticket           `json:"lskey,omitempty"`
	SKey     *Ticket           `json:"skey,omitempty"`
	SID      *Ticket           `json:"sid,omitempty"`
	Sig64    *Ticket           `json:"sig64,omitempty"`
	SuperKey *Ticket           `json:"super_key,omitempty"`
	ST       *Ticket           `json:"st,omitempty"`
	STWeb    *Ticket           `json:"stweb,omitempty"`
	VKey     *Ticket           `json:"vkey,omitempty"`
	Domains  map[string]string `json:"domains,omitempty"`
	KSID     []byte            `json:"ksid,omitempty"`
}

type Ticket struct {
	Key Key16Bytes `json:"key,omitempty"`
	Sig []byte     `json:"sig,omitempty"`
	Exp int64      `json:"exp,omitempty"`
	Iss int64      `json:"iss,omitempty"`
}

func (t Ticket) Valid() bool {
	return time.Now().Before(time.Unix(t.Exp, 0))
}

type Key16Bytes [16]byte

func (v Key16Bytes) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(hex.EncodeToString(v[:]))), nil
}

func (v *Key16Bytes) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	p, err := hex.DecodeString(unquoted)
	if err != nil {
		return err
	}
	if v == nil {
		v = new(Key16Bytes)
	}
	copy(v[:], p[:])
	return nil
}

func (v *Key16Bytes) Get() [16]byte {
	if v == nil {
		v = new(Key16Bytes)
	}
	return *v
}

func (v *Key16Bytes) Set(b [16]byte) {
	if v == nil {
		v = new(Key16Bytes)
	}
	copy(v[:], b[:])
}

type FakeSource struct {
	App    *FakeApp
	Device *FakeDevice
}

type FakeApp struct {
	FixID int32
	AppID int32

	PkgName  string
	VerName  string
	Revision string
	SigHash  [16]byte

	BuildAt int64
	SDKVer  string
	SSOVer  uint32

	ImageType  uint8
	MiscBitMap uint32

	CanCaptcha bool
}

type FakeDevice struct {
	OS FakeDeviceOS

	APNName   []byte
	SIMOPName []byte

	Bootloader   string
	ProcVersion  string
	Codename     string
	Incremental  string
	Fingerprint  string
	BootID       string
	Baseband     string
	InnerVersion string

	NetworkType  uint8 // 0x00: Others; 0x01: Wi-Fi
	NetIPFamily  uint8 // 0x00: Others; 0x01: IPv4; 0x02: IPv6; 0x03: Dual
	IPv4Address  net.IP
	IPv6Address  net.IP
	MACAddress   string
	BSSIDAddress string
	SSIDAddress  string

	IMEI string
	IMSI string
	GUID [16]byte // []byte("%4;7t>;28<fclient.5*6")

	GUIDFlag      uint32
	IsGUIDFileNil bool
	IsGUIDGenSucc bool
	IsGUIDChanged bool
}

type FakeDeviceOS struct {
	Type        string
	Version     string
	BuildBrand  []byte
	BuildID     string
	BuildModel  string
	SDKVersion  uint32
	NetworkType uint16 // 0x0002: Wi-Fi
}
