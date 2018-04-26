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

//go:generate protoc --go_out=plugins=grpc:. pb/protocol.proto
//go:generate abigen -pkg trans -sol trans/contracts/sink.sol -out trans/sink.sol.go
//go:generate abigen -pkg trans -sol trans/contracts/wallet.sol -out trans/wallet.sol.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"
	"github.com/boxproject/voucher/common"
	"github.com/boxproject/voucher/config"
	//"github.com/boxproject/voucher/http"
	logger "github.com/alecthomas/log4go"
	"github.com/boxproject/voucher/localdb"
	"github.com/boxproject/voucher/node"
	"github.com/boxproject/voucher/operate"
	"github.com/boxproject/voucher/token"
	"github.com/boxproject/voucher/trans"
	"github.com/boxproject/voucher/util"
	"github.com/mdp/qrterminal"
	"gopkg.in/urfave/cli.v1"
	"github.com/awnumar/memguard"
)

func main() {
	app := newApp()
	if err := app.Run(os.Args); err != nil {
		logger.Error("start app failed. %v", err)
	}
}

func run(ctx *cli.Context) (err error) {
	defer  memguard.DestroyAll()
	var (
		cfg *config.Config
		db  localdb.Database
	)
	logger.LoadConfiguration("log.xml")

	//config
	filePath := ctx.String("c")
	if cfg, err = common.LoadConfig(filePath); err != nil {
		logger.Error("load config failed. cause: %v", err)
		return err
	}

	printQrCode(cfg)

	//log
	//if err = log.InitLogger(&cfg.LogConfig); err != nil {
	//	log.Errorf("init logger failed. cause: %v", err)
	//	return err
	//}

	//defer log.Close()

	logger.Debug("config file:\n%s", cfg)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh,
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP, syscall.SIGKILL,
		syscall.SIGUSR1, syscall.SIGUSR2)

	//db
	if db, err = localdb.RunDatabase(cfg); err != nil {
		fmt.Println("run database failed. cause: %v", err)
		//logger.Error("run database failed. cause: %v", err)
		return err
	}
	defer db.Close()

	//init status
	initStatus(cfg, db)

	//init eth token
	initEthToken(db)
	//grpc
	nodeApi, err := node.InitConn(cfg, db)
	if err != nil {
		fmt.Println("rpc client init failed. cause: %v", err)
		return err
	}
	defer nodeApi.CloseConn()
	/*
		//down web api
		downloader, err := http.NewCertService(cfg, db)
		if err != nil {
			log.Errorf("new cert download service failed. cause: %v", err)
			return err
		}

		if err = downloader.Start(); err != nil {
			log.Errorf("start cert download service failed. cause: %v", err)
			return
		}
	*/
	//eth conn
	ethHander, err := trans.NewEthHandler(cfg, db)
	defer ethHander.Stop()

	if err != nil {
		fmt.Println("new default handler config failed. cause: %v", err)
		return err
	}
	oHandler := operate.InitOperateHandler(cfg, db, ethHander)
	go oHandler.Start()
	defer oHandler.Close()

	////web api
	//cliapi, err := http.NewClientAPI(cfg, db, nodeApi, ethHander)
	//if err != nil {
	//	logger.Error("new client api service failed. cause: %v", err)
	//	return err
	//}
	//if err = cliapi.Start(); err != nil {
	//	logger.Error("start cert cliapi service failed. cause: %v", err)
	//	return err
	//}

	//account rep

	//defer cliapi.Shutdown()

	//printQrCode(cfg)

	<-quitCh
	//downloader.Shutdown(context.Background())

	return err
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Action = run

	app.Version = PrintVersion(gitCommit, stage, version)
	app.Name = "voucher"
	app.Usage = "command line interface"
	app.Author = "BOX.la"
	app.Copyright = "Copyright 2017-2019 The BOX.la Authors"
	app.Email = "develop@2se.com"
	app.Description = "The automatic teller machine for cryptocurrency"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c,config",
			Usage: "Path of the config.toml file",
			Value: "",
		},
	}

	return app
}

//init eth tokens
func initEthToken(db localdb.Database) {
	token.LoadTokenFrom(db)
}

//init status
func initStatus(cfg *config.Config, db localdb.Database) {
	if statusBytes, err := db.Get([]byte(config.STATAUS_KEY)); err != nil {
		logger.Info("load status from db err:%s", err)
		config.RealTimeStatus.ServerStatus = config.VOUCHER_STATUS_UNCREATED //TODO
		config.RealTimeStatus.Total = cfg.Secret.AppNum

		config.RealTimeStatus.CoinStatus = append(config.RealTimeStatus.CoinStatus, config.CoinStatu{Name: config.COIN_NAME_BTC, Category: config.CATEGORY_BTC, Decimals: config.COIN_DECIMALS_BTC, Used: false})
		//config.RealTimeStatus.CoinStatus = append(config.RealTimeStatus.CoinStatus, config.CoinStatu{Name: config.COIN_NAME_ETH, Category: config.CATEGORY_ETH, Decimals: config.COIN_DECIMALS_ETH, Used: true})
		//config.RealTimeStatus.KeyStroeStatus = make([]config.KeyStroeStatu)
	} else {
		if err = json.Unmarshal(statusBytes, config.RealTimeStatus); err != nil {
			logger.Error("unmarshal status err: %s", err)
		} else {
			if config.RealTimeStatus.ServerStatus == config.VOUCHER_STATUS_STATED {
				config.RealTimeStatus.ServerStatus = config.VOUCHER_STATUS_PAUSED
			}
			if config.RealTimeStatus.CoinStatus == nil || len(config.RealTimeStatus.CoinStatus) == 0 {
				config.RealTimeStatus.CoinStatus = append(config.RealTimeStatus.CoinStatus, config.CoinStatu{Name: config.COIN_NAME_BTC, Category: config.CATEGORY_BTC, Decimals: config.COIN_DECIMALS_BTC, Used: false})
			}
			config.RealTimeStatus.Status = config.PASSWORD_STATUS_OK
			config.RealTimeStatus.NodesAuthorized = nil
		}
	}
}

func PrintVersion(gitCommit, stage, version string) string {
	if gitCommit != "" {
		return fmt.Sprintf("%s-%s-%s", stage, version, gitCommit)
	}
	return fmt.Sprintf("%s-%s", stage, version)
}

func printQrCode(cfg *config.Config) {
	qrCodeArrays := []string{cfg.AgentSerCfg.IpAndPort, util.GetAesKeyRandom()}
	qrCodeByte, _ := json.Marshal(qrCodeArrays)
	fmt.Println("keys:", string(qrCodeByte))
	qrterminal.Generate(string(qrCodeByte), qrterminal.L, os.Stdout)
}
