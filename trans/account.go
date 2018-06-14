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

package trans

import (
	"bytes"
	"crypto/ecdsa"
	"errors"

	"crypto/md5"
	log "github.com/alecthomas/log4go"
	"github.com/awnumar/memguard"
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
	verror "github.com/boxproject/voucher/errors"
	"github.com/boxproject/voucher/localdb"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"sort"
	"unsafe"
)

var pk *memguard.LockedBuffer

func GetSecret(db localdb.Database) (string, error) {
	v, err := db.Get(config.SECRET.Bytes())
	if err != nil {
		log.Error("get secret from db failed. casue: %v", err)
		return "", err
	}
	//d := binary.LittleEndian.Uint32(v)
	return string(v), nil
}

func RecoverPrivateKey(db localdb.Database, args ...[]byte) error {
	sort.Sort(IntSlice(args))
	v, err := db.Get(config.SECRET.Bytes())
	if err != nil {
		log.Error("get secret from db failed. casue: %v", err)
		return err
	}

	d := sha3.NewKeccak512()
	d.Write(v)
	for _, b := range args {
		d.Write(b)
	}
	dBytes := d.Sum(nil)
	hahsByte := md5.Sum(dBytes)
	db.Put(config.PRIVATEKEYHASH.Bytes(), hahsByte[:]) //private key hash

	buf := bytes.NewBuffer(dBytes)

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), buf)
	if err != nil {
		log.Error("generate key error:%v", err)
		return err
	}
	pk, err = memguard.NewImmutableFromBytes(crypto.FromECDSA(privateKeyECDSA))
	if err != nil {
		return err
	}

	if err := setPubKeytoDB(db); err != nil {
		log.Error("generate pubkey error:%v", err)
		return err
	}
	return nil
}

//校验私钥hash
//Verify private key hash
func CheckPrivateKey(db localdb.Database, args ...[]byte) bool {
	sort.Sort(IntSlice(args))
	if pvhash, err := db.Get(config.PRIVATEKEYHASH.Bytes()); err != nil {
		log.Error("CheckPrivateKey err: %s", err)
		return false
	} else {
		v, err := db.Get(config.SECRET.Bytes())
		if err != nil {
			log.Error("get secret from db failed. casue: %v", err)
			return false
		}

		d := sha3.NewKeccak512()
		d.Write(v)
		sort.Sort(IntSlice(args))
		for _, b := range args {
			d.Write(b)
		}
		pv := d.Sum(nil)
		hashBytes := md5.Sum(pv)
		if bytes.Equal(pvhash, hashBytes[:]) {
			log.Debug("pvhash is equal dbHash")
			buf := bytes.NewBuffer(pv)
			privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), buf)
			if err != nil {
				log.Error("generate key error:%v", err)
				return false
			}
			pk, err = memguard.NewImmutableFromBytes(crypto.FromECDSA(privateKeyECDSA))
			if err != nil {
				return false
			}
			if err := setPubKeytoDB(db); err != nil {
				log.Error("generate pubkey error:%v", err)
				return false
			}

			return true
		}
	}
	return false
}

func setPubKeytoDB(db localdb.Database) error {
	if pk == nil {
		return errors.New("get eth privkey error, data is nil")
	}
	keyArrayPtr := (*[32]byte)(unsafe.Pointer(&pk.Buffer()[0]))
	privateKeyECDSA, err := crypto.ToECDSA(keyArrayPtr[:])
	if err != nil {
		return err
	}
	ethPubKey := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()
	if err := db.Put([]byte(config.PUBKEY_ETH), []byte(ethPubKey)); err != nil {
		return err
	}

	pv := btcec.PublicKey(privateKeyECDSA.PublicKey)
	btcPubKey, _ := btcutil.NewAddressPubKey(pv.SerializeCompressed(), &chaincfg.MainNetParams)
	if err := db.Put([]byte(config.PUBKEY_BTC), []byte(btcPubKey.EncodeAddress())); err != nil {
		return err
	}

	return nil
}

//是否已创建过key
//If the key has been created
func ExistPrivateKeyHash(db localdb.Database) bool {
	data, err := db.Get(config.PRIVATEKEYHASH.Bytes())
	if err == verror.NoDataErr {
		log.Info("ExistPrivateKeyHash no data")
	} else {
		log.Error("ExistPrivateKeyHash err: %s", err)
	}
	if len(data) > 0 {
		return true
	}
	return false
}

func NewKeyedTransactor() *bind.TransactOpts {
	if pk == nil {
		log.Error("get eth privkey error, data is nil")
		return nil
	} else {
		keyArrayPtr := (*[32]byte)(unsafe.Pointer(&pk.Buffer()[0]))
		privateKeyECDSA, err := crypto.ToECDSA(keyArrayPtr[:])
		if err != nil {
			log.Error(err)
			return nil
		}
		keyAddr := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
		return &bind.TransactOpts{
			From: keyAddr,
			Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if address != keyAddr {
					return nil, errors.New("not authorized to sign this account")
				}
				signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privateKeyECDSA)
				if err != nil {
					return nil, err
				}
				return tx.WithSignature(signer, signature)
			},
			//GasLimit: uint64(config.DefPubEthGasLimit),
		}
	}
}

func GetBtcPrivKey() txscript.KeyDB {
	if pk == nil {
		err := errors.New("get btc privkey error, data is nil")
		log.Error(err)
		return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
			return nil, true, err
		})
	} else {
		return txscript.KeyClosure(func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
			keyArrayPtr := (*[32]byte)(unsafe.Pointer(&pk.Buffer()[0]))
			privateKeyECDSA, err := crypto.ToECDSA(keyArrayPtr[:])
			if err != nil {
				log.Error("crypto.ToECDSA err:", err)
				return nil, true, err
			} else {
				pv := btcec.PrivateKey(*privateKeyECDSA)
				return &pv, true, nil
			}

		})
	}
}

func GetBtcPubKey(db localdb.Database) (string, error) {
	if datas, err := db.Get([]byte(config.PUBKEY_BTC)); err == nil {
		return string(datas), nil
	} else {
		return "", err
	}
}

func PrivateKeyToHex(db localdb.Database) string {
	if datas, err := db.Get([]byte(config.PUBKEY_ETH)); err == nil {
		return string(datas)
	} else {
		return ""
	}
}

type IntSlice [][]byte

func (s IntSlice) Len() int           { return len(s) }
func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IntSlice) Less(i, j int) bool { return string(s[i]) < string(s[j]) }
