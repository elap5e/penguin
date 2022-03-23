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

package ecdh

import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

var X509Prefix, _ = hex.DecodeString("3059301306072a8648ce3d020106082a8648ce3d030107034200")

type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

func (pub *PublicKey) Bytes() []byte {
	return append(append([]byte{0x04}, pub.X.Bytes()...), pub.Y.Bytes()...)
}

type PrivateKey struct {
	PublicKey
	D *big.Int
}

func (priv *PrivateKey) Public() *PublicKey {
	return &priv.PublicKey
}

func (priv *PrivateKey) SharedSecret(pub *PublicKey) [16]byte {
	sx, _ := priv.PublicKey.Curve.ScalarMult(pub.X, pub.Y, priv.D.Bytes())
	return md5.Sum(sx.Bytes()[:16])
}

func GenerateKey() (*PrivateKey, error) {
	c := elliptic.P256()
	k, x, y, err := elliptic.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, err
	}
	priv := &PrivateKey{}
	priv.Curve = c
	priv.PublicKey.X, priv.PublicKey.Y = x, y
	priv.D = new(big.Int).SetBytes(k)
	return priv, nil
}
