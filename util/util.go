package util

import (
	//"bytes"
	log "github.com/alecthomas/log4go"
	"io/ioutil"
	//"math/big"
	"bytes"
	"crypto/aes"
	//"crypto/rand"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/localdb"
	"math/big"
	"math/rand"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
	"errors"
)

var (
	rootPath string
	filePath string
	//aesRandomStr string
)

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}

	return ""
}

func GetFilePath() string {
	return filePath
}

func DefaultConfigDir() string {
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".bcmonitor")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "bcmonitor")
		} else {
			return filepath.Join(home, ".bcmonitor")
		}
	}

	return ""
}

// configPath 不为空时，不检查fileName
func GetConfigFilePath(configPath, defaultFileName string) string {
	for i := 0; i < 3; i++ {
		if configPath != "" {
			if _, err := os.Stat(configPath); !os.IsNotExist(err) {
				break
			}
		}

		if i == 0 {
			configPath = path.Join(GetFilePath(), defaultFileName)
		} else if i == 1 {
			configPath = path.Join(DefaultConfigDir(), defaultFileName)
		}
	}

	return configPath
}

//file
var NoRWMutex sync.RWMutex

func WriteNumberToFile(filePath string, blkNumber *big.Int) error {
	return ioutil.WriteFile(filePath, []byte(blkNumber.String()), 0755)
}

func ReadNumberFromFile(filePath string) (*big.Int, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Debug("file not found, %v", err)
		return big.NewInt(0), err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	data = bytes.TrimSpace(data)
	delta, isok := big.NewInt(0).SetString(string(data), 10)
	if isok == false {
		return big.NewInt(0), errors.New("data is nil")
	}
	return delta, nil
}

//func ReadBlockInfoFromFile(filePath string) (*config.BlockInfoCfg, error) {
//	noRWMutex.Lock()
//	defer noRWMutex.Unlock()
//	blkInfoCfg := &config.BlockInfoCfg{big.NewInt(0),big.NewInt(0),big.NewInt(0)}
//	if _, err := os.Stat(filePath); os.IsNotExist(err) {
//		log.Debug("file not found, %v", err)
//		return blkInfoCfg, nil
//	}
//	data, err := ioutil.ReadFile(filePath)
//	if err != nil {
//		return blkInfoCfg, err
//	}
//	if err = json.Unmarshal(data, blkInfoCfg); err != nil {
//		log.Debug("json unmarshal error:%s", err)
//	}
//	return blkInfoCfg, nil
//}
//
//func WriteBlockInfoToFile(filePath string, blkInfoCfg *config.BlockInfoCfg) error {
//	noRWMutex.Lock()
//	defer noRWMutex.Unlock()
//	if data, err := json.Marshal(blkInfoCfg); err != nil {
//		log.Debug("json unmarshal error:%s", err)
//	} else {
//		ioutil.WriteFile(filePath, data, 0755)
//	}
//	return nil
//}

func GetCurrentIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		log.Error("Get local IP addr failed!!!")
		return "localhost"
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

//var cbcKey string//default
//
//func GetAesKeyRandom() string {
//	if cbcKey == "" {
//		r := rand.New(rand.NewSource(time.Now().UnixNano())) //添加时间生成随机数
//		key := make([]byte, aes.BlockSize)
//		copy(key, []byte(string(r.Intn(100000))))
//		cbcKey = string(key)
//		//更新key
//	}
//	return cbcKey
//}

//aes key 存入db
func GetAesKeyRandomFromDb(db localdb.Database) []byte {
	if aesKeyBytes, err := db.Get([]byte(config.APP_KEY_PRIFIX)); err != nil {
		log.Info("get app key failed. err : %s", err)
		cbcKey := make([]byte, aes.BlockSize)

		r := rand.New(rand.NewSource(time.Now().UnixNano())) //添加时间生成随机数
		ri := r.Uint64()
		rstr := strconv.FormatUint(ri, 16)
		copy(cbcKey, []byte(rstr))
		if err = db.Put([]byte(config.APP_KEY_PRIFIX), cbcKey); err != nil {
			log.Debug("land aes to db err: %s", err)
			return config.DefAesKey
		}
		return cbcKey
	} else {
		return aesKeyBytes
	}

}
