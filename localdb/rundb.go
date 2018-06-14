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
package localdb

import (
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/errors"
)

func RunDatabase(cfg *config.Config) (Database, error) {
	db, err := NewDatabase(&cfg.DB)
	if err != nil {
		return nil, err
	}

	_, err = db.Get(config.INIT.Bytes())
	if err != nil {
		if errors.NoDataErr != err {
			return nil, err
		}

		if err = prepare(db, cfg); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func prepare(db Database, cfg *config.Config) (err error) {
	// 如果没有初始化过，则需要产生证书
	// Create certificate if no initialization yet
	if err = common.GenCert(cfg, config.SVRCERTSCRIPT); err != nil {
		errors.GenServerCertErr.Err = err
		return errors.GenServerCertErr
	}

	if err = common.GenCert(cfg, config.CLTCERTSCRIPT); err != nil {
		errors.GenClientCertErr.Err = err
		return errors.GenClientCertErr
	}

	var (
		secret []byte
		pass   []byte
	)

	// 产生服务器端密钥
	// Generate server side private key
	if secret, err = common.GenSecret(cfg.Secret.SecretLength); err != nil {
		errors.GenServerSecretErr.Err = err
		return errors.GenServerSecretErr.Err
	}

	// 产生随机密钥
	// Generate random private key
	if pass, err = common.GenSecret(cfg.Secret.PassLength); err != nil {
		errors.GenRandomPassErr.Err = err
		return errors.GenRandomPassErr
	}

	if err = db.Put(config.SECRET.Bytes(), secret); err != nil {
		errors.SaveSecretErr.Err = err
		return errors.SaveSecretErr
	}

	if err = db.Put(config.PASS.Bytes(), pass); err != nil {
		errors.SavePassErr.Err = err
		return errors.SavePassErr
	}

	if err = db.Put(config.INIT.Bytes(), config.Inited.Bytes()); err != nil {
		return err
	}

	return nil
}
