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
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	x509Prefix, _ = hex.DecodeString("3059301306072a8648ce3d020106082a8648ce3d030107034200")
)

var (
	serverECDHPublicKey, _ = hex.DecodeString("04ebca94d733e399b2db96eacdd3f69a8bb0f74224e2b44e3357812211d2e62efbc91bb553098e25e33a799adc7f76feb208da7c6522cdb0719a305180cc54a82e")
	serverRSAPublicKey, _  = base64.StdEncoding.DecodeString("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuJTW4abQJXeVdAODw1CamZH4QJZChyT08ribet1Gp0wpSabIgyKFZAOxeArcCbknKyBrRY3FFI9HgY1AyItH8DOUe6ajDEb6c+vrgjgeCiOiCVyum4lI5Fmp38iHKH14xap6xGaXcBccdOZNzGT82sPDM2Oc6QYSZpfs8EO7TYT7KSB2gaHz99RQ4A/Lel1Vw0krk+DescN6TgRCaXjSGn268jD7lOO23x5JS1mavsUJtOZpXkK9GqCGSTCTbCwZhI33CpwdQ2EHLhiP5RaXZCio6lksu+d8sKTWU1eEiEb3cQ7nuZXLYH7leeYFoPtbFV4RicIWp0/YG+RP7rLPCwIDAQAB")
)

var (
	ServerKeyVersion = int16(0)
	ServerPublicKey  = &PublicKey{}
)

func init() {
	setServerPublicKey(serverECDHPublicKey, 1)
	updateServerPublicKey()
}

func setServerPublicKey(key []byte, ver int16) error {
	pub, err := x509.ParsePKIXPublicKey(append(x509Prefix, key...))
	if err != nil {
		return err
	}
	ServerKeyVersion = ver
	ServerPublicKey.Curve = pub.(*ecdsa.PublicKey).Curve
	ServerPublicKey.X = pub.(*ecdsa.PublicKey).X
	ServerPublicKey.Y = pub.(*ecdsa.PublicKey).Y
	return nil
}

func updateServerPublicKey() error {
	type serverPublicKey struct {
		QuerySpan         int32 `json:"QuerySpan"`
		PublicKeyMetadata struct {
			KeyVersion    int16  `json:"KeyVer"`
			PublicKey     string `json:"PubKey"`
			PublicKeySign string `json:"PubKeySign"`
		} `json:"PubKeyMeta"`
	}
	req, err := http.NewRequest("GET", "https://keyrotate.qq.com/rotate_key?cipher_suite_ver=305&uin=10000", nil)
	if err != nil {
		return err
	}
	// req.Header.Add("js.fetch:mode", "no-cors")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := serverPublicKey{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	pub, err := x509.ParsePKIXPublicKey(serverRSAPublicKey)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(fmt.Sprintf("305%d%s", data.PublicKeyMetadata.KeyVersion, data.PublicKeyMetadata.PublicKey)))
	sig, _ := base64.StdEncoding.DecodeString(
		data.PublicKeyMetadata.PublicKeySign,
	)
	if err := rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hashed[:], sig); err != nil {
		return err
	}
	key, _ := hex.DecodeString(data.PublicKeyMetadata.PublicKey)
	if err := setServerPublicKey(key, data.PublicKeyMetadata.KeyVersion); err != nil {
		return err
	}
	return nil
}
