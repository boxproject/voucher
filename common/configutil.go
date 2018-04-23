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
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/errors"
	"github.com/BurntSushi/toml"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func LoadConfig(filePath string) (*config.Config, error) {
	var err error
	if filePath == "" {
		if filePath, err = tryFindConfig(); err != nil {
			return nil, err
		}
	}

	var c config.Config
	if _, err := toml.DecodeFile(filePath, &c); err != nil {
		return nil, err
	}

	if c.Basedir == "" {
		return nil, errors.NoBasedirErr
	}

	if len(c.Node) == 0 {
		return nil, errors.NoPrivateNodeErr
	}

	if c.EthereumConfig.GasLimit > 0 {//修改默认gaslimit
		config.DefPubEthGasLimit = c.EthereumConfig.GasLimit
	}

	return setDefault(&c), nil
}

func setDefault(cfg *config.Config) *config.Config {
	if cfg.Secret.SecretLength <= 0 {
		cfg.Secret.SecretLength = config.DEFAULTSECRETLEN
	}

	if cfg.Secret.PassLength <= 0 {
		cfg.Secret.PassLength = config.DEFAULTPASSLENT
	}

	if cfg.Secret.AppNum <= 0 {
		cfg.Secret.AppNum = config.DEFAULTAPPNUM
	}

	if cfg.DB.Filepath == "" {
		cfg.DB.Filepath = path.Join(cfg.Basedir, config.DEFAULTDBDIRNAME)
	}

	if cfg.DB.Prefix == "" {
		cfg.DB.Prefix = config.DEFAULTKEYPREFIX
	}

	if cfg.LogConfig.LogPath == "" {
		cfg.LogConfig.LogPath = path.Join(cfg.Basedir, config.DEFAULTLOGDIRNAME)
	}

	cfg.LogConfig.Encoding = strings.ToLower(cfg.LogConfig.Encoding)

	return cfg
}

func tryFindConfig() (string, error) {
	self, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	p, err := filepath.Abs(self)
	if err != nil {
		return "", err
	}

	filePath := path.Join(path.Dir(p), config.DEFAULTCONFIGFILENAME)

	_, err = os.Stat(filePath)
	if err == nil {
		return filePath, nil
	}

	if os.IsNotExist(err) {
		return "", errors.NoConfigErr
	}

	return "", err
}
