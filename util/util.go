package util

import (
	//"bytes"
	"fmt"
	log "github.com/alecthomas/log4go"
	"io/ioutil"
	//"math/big"
	"bytes"
	"crypto/aes"
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
	"sync"
	"time"
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
		fmt.Println("i......", i)
		if configPath != "" {
			if _, err := os.Stat(configPath); !os.IsNotExist(err) {
				fmt.Println("Stat....", err)
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
var noRWMutex sync.RWMutex

func WriteNumberToFile(filePath string, blkNumber *big.Int) error {
	noRWMutex.Lock()
	defer noRWMutex.Unlock()
	return ioutil.WriteFile(filePath, []byte(blkNumber.String()), 0755)
}

func ReadNumberFromFile(filePath string) (*big.Int, error) {
	noRWMutex.Lock()
	defer noRWMutex.Unlock()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Debug("file not found, %v", err)
		return big.NewInt(0), nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	data = bytes.TrimSpace(data)
	delta, _ := big.NewInt(0).SetString(string(data), 10)
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

// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func GetIntnRandom(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //添加时间生成随机数
	return r.Intn(n)
}

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

var cbcKey string = "abcdefghijklmnop" //default

func GetAesKeyRandom() string {
	if cbcKey == "" {
		r := rand.New(rand.NewSource(time.Now().UnixNano())) //添加时间生成随机数
		key := make([]byte, aes.BlockSize)
		copy(key, []byte(string(r.Intn(100000))))
		cbcKey = string(key)
		//更新key
	}
	return cbcKey
}

//aes key 存入db
func GetAesKeyRandomFromDb(db localdb.Database) string {
	if aesKeyBytes, err := db.Get([]byte(config.APP_KEY_PRIFIX)); err != nil {
		log.Info("get app key failed. err : %s", err)
		r := rand.New(rand.NewSource(time.Now().UnixNano())) //添加时间生成随机数
		key := make([]byte, aes.BlockSize)
		copy(key, []byte(string(r.Intn(10000000))))
		cbcKey = string(key)
		if err = db.Put([]byte(config.APP_KEY_PRIFIX), key); err != nil {
			log.Debug("land aes to db err: %s", err)
		}
	} else {
		cbcKey = string(aesKeyBytes)
	}
	return cbcKey
}
