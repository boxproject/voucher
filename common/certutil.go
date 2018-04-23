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
	"crypto/rand"
	"io"
	"os/exec"
	"path"

	"github.com/boxproject/voucher/config"
)

func GenCert(cfg *config.Config, certScriptFile string) error {
	certPath := path.Join(cfg.Basedir, config.DEFAULTCERTSDIRNAME)
	script := path.Join(cfg.Basedir, config.DEFAULTSCRIPTDIRNAME, certScriptFile)
	cmd := exec.Command(script, certPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func GenSecret(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GetServerCert(cfg *config.Config) string {
	return path.Join(cfg.Basedir, config.DEFAULTCERTSDIRNAME, "server.pem")
}

func GetServerKey(cfg *config.Config) string {
	return path.Join(cfg.Basedir, config.DEFAULTCERTSDIRNAME, "server.key")
}
