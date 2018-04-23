// Copyright 2017. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type (
	Hash     = common.Hash
	Address  = common.Address
	EtherLog = types.Log
)
const AddressLength = common.AddressLength


func BytesToHex(data []byte, with0x bool) string {
	unchecksummed := hex.EncodeToString(data)
	result := []byte(unchecksummed)
	hash := crypto.Keccak256(result)

	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}

	if with0x {
		return "0x" + string(result)
	}

	return string(result)
}

func BytesToAddress(data []byte) common.Address {
	return common.BytesToAddress(data)
}

func HexToAddress(data string) common.Address {
	return common.HexToAddress(data)
}

func BytesToHash(data []byte) common.Hash {
	return common.BytesToHash(data)
}

func HexToBytes(hexstr string) ([]byte, error) {
	return hex.DecodeString(hexstr)
}

func HexToHash(hexstr string) common.Hash {
	return common.HexToHash(hexstr)
}

func SignEvent(f string) Hash {
	data := []byte(f)
	d := sha3.NewKeccak256()
	d.Write(data)
	return common.BytesToHash(d.Sum(nil))
}

func AddressEquals(a, b Address) bool {
	return bytes.Equal(a.Bytes(), b.Bytes())
}

func Byte2Byte32(src []byte) [32]byte {
	var obj [32]byte
	copy(obj[:], src[:32])
	return obj
}
