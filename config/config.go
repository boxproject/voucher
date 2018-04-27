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
package config

import (
	"fmt"
)

const (
	SVRCERTSCRIPT = "server-cert.sh"
	CLTCERTSCRIPT = "client-cert.sh"

	DEFAULTSECRETLEN      = 40
	DEFAULTPASSLENT       = 8
	DEFAULTAPPNUM         = 3
	DEFAULTKEYPREFIX      = "BOX"
	DEFAULTDBDIRNAME      = "db"
	DEFAULTCERTSDIRNAME   = "certs"
	DEFAULTSCRIPTDIRNAME  = "scripts"
	DEFAULTLOGDIRNAME     = "log/voucher.log"
	DEFAULTCONFIGFILENAME = "config.toml"
)

type Config struct {
	Basedir        string             `toml:"basedir"`
	AgentSerCfg    AgentServiceConfig `toml:"agentservice"`
	Secret         SecretConfig       `toml:"secret"`
	DB             DatabaseConfig     `toml:"database"`
	LogConfig      LogConfig          `toml:"log"`
	CertHTTPConfig CertHTTPConfig     `toml:"service"`
	APIConfig      APIConfig          `toml:"api"`
	EthereumConfig EthereumConfig     `toml:"ethereum"`
	BitcoinConfig  BitcoinConfig      `toml:"bitcoin"` //add by john.yang
}

func (c *Config) String() string {
	return fmt.Sprintf("[Base DIR] %s\n[Database]\n%s\n\n[Nodes]\n%s\n", c.Basedir, c.DB)
}

type AgentServiceConfig struct {
	Name           string `toml:"name"`
	Alias          string `toml:"alias"`
	CompanionName  string `toml:"companionname"`
	IpAndPort      string `toml:"ipandport"`
	Pattern        string `toml:"pattern"`
	RPCAddr        string `toml:"rpcapi"`
	ClientCertPath string `toml:"clientcert"`
	ClientKeyPath  string `toml:"clientkey"`
}

type DatabaseConfig struct {
	Filepath  string `toml:"filePath"`
	Cache     int    `toml:"cache"`
	Openfiles int    `toml:"openFiles"`
	Prefix    string `toml:"prefix"`
}

func (d DatabaseConfig) String() string {
	return fmt.Sprintf("[File Path]: %s, [Cache] %d, [Openfiles] %d", d.Filepath, d.Cache, d.Openfiles)
}

type SecretConfig struct {
	SecretLength int `toml:"secretLength"`
	PassLength   int `toml:"passLength"`
	AppNum       int `toml:"appNum"`
}

// NodeConfig - 私链连接配置
type NodeConfig struct {
	Name           string `toml:"name"`
	RPCAddr        string `toml:"rpcapi"`
	ClientCertPath string `toml:"clientcert"`
	ClientKeyPath  string `toml:"clientkey"`
}

func (n NodeConfig) String() string {
	return fmt.Sprintf("[%s]: RPC API: %s\nClient cert: %s\nClient Key: %s\n", n.Name, n.RPCAddr, n.ClientCertPath, n.ClientKeyPath)
}

type LogConfig struct {
	Level   string `toml:"level"`
	LogPath string `toml:"logPath"`
	// json/console
	Encoding string `toml:"encoding"`
}

type CertHTTPConfig struct {
	APIConfig
}

type APIConfig struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}

type EthereumConfig struct {
	Scheme          string `toml:"scheme"`
	DelayedBlocks   int64  `toml:"delayedBlocks"`
	CursorBlocks    int64  `toml:"cursorBlocks"`
	Retries         int    `toml:retries`
	GasLimit        int    `toml:gasLimit`
	AccountPoolSize int    `toml:accountPoolSize`
	BlockNoFilePath string `toml:blockNoFilePath`
	NonceFilePath   string `toml:nonceFilePath`
	CursorFilePath  string `toml:cursorFilePath`
}

type BitcoinConfig struct {
	Type          	 string   `toml:"type"`
	Host            string `toml:"host"`
	Rpcuser         string `toml:"rpcuser"`
	Rpcpass         string `toml:"rpcpass"`
	Clientcert      string `toml:"clientcert"`
	Clientkey       string `toml:"clientkey"`
	Confirmations   int64 `toml:"confirmations"`
	Initusernum     int    `toml:initusernum`
	Initheight      int64  `toml:initheight`
	BlockNoFilePath string `toml:blockNoFilePath`
}
