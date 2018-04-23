package util

import (
	//"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/boxproject/voucher/config"
	"github.com/boxproject/voucher/localdb"
	//"encoding/base64"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	log "github.com/alecthomas/log4go"
)

//var decrypted string
//var privateKey, publicKey []byte

//func initi() {
//	var err error
//	flag.StringVar(&decrypted, "d", "", "加密过的数据")
//	flag.Parse()
//	publicKey, err = ioutil.ReadFile("public.pem")
//	if err != nil {
//		os.Exit(-1)
//	}
//	privateKey, err = ioutil.ReadFile("private.pem")
//	if err != nil {
//		os.Exit(-1)
//	}
//}

//func main() {
//	var data []byte
//	var err error
//	data, err = RsaEncrypt([]byte("fyxichen"))
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("RsaEncrypt...", data)
//	origData, err := RsaDecrypt(data)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(origData))
//}

//func getPublicKey(db localdb.Database, appId string) {
//	if keyBytes, err := db.Get([]byte(config.APP_KEY_PRIFIX + appId)); err != nil {
//		publicKey = keyBytes
//	}
//}

//// 加密
//func RsaEncrypt(origData []byte) ([]byte, error) {
//	block, _ := pem.Decode(publicKey)
//	if block == nil {
//		return nil, errors.New("public key error")
//	}
//	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//	pub := pubInterface.(*rsa.PublicKey)
//	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
//}

// 解密
//func RsaDecrypt(ciphertext []byte) ([]byte, error) {
//	block, _ := pem.Decode(privateKey)
//	if block == nil {
//		return nil, errors.New("private key error!")
//	}
//	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
//}

//公钥验证
func RsaSignVer(data []byte, signature []byte, db localdb.Database, appId string) error {
	if publicKey, err := db.Get([]byte(config.APP_KEY_PRIFIX + appId)); err != nil {
		return err
	} else {
		deSignData, err := base64.StdEncoding.DecodeString(string(signature))
		if err != nil {
			log.Error("base64 Decode error:", err)
		}
		hashed := sha256.Sum256(data)
		keybase64, err := base64.RawStdEncoding.DecodeString(string(publicKey))
		if err != nil {
			log.Error("base64 Decode error:", err)
		}
		pub, err := x509.ParsePKCS1PublicKey(keybase64)
		if err != nil {
			log.Error("ParsePKCS1PublicKey err: %s", err)
			return err
		}

		//block, _ := pem.Decode(publicKey)
		//if block == nil {
		//	return errors.New("public key error")
		//}
		//// 解析公钥
		//pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		//if err != nil {
		//	log.Error("解析公钥")
		//	return err
		//}
		//// 类型断言
		//pub := pubInterface.(*rsa.PublicKey)
		// 验证签名
		return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], deSignData)
	}
}

//公钥验证with公钥
func RsaSignVerWithKey(data []byte, signature []byte, publicKey []byte) error {
	log.Debug("publicKey.....", string(publicKey))
	bSignData, err := base64.StdEncoding.DecodeString(string(signature))
	if err != nil {
		log.Error("base64 Decode error:", err)
		return err
	}
	hashed := sha256.Sum256(data)
	bKey, err := base64.RawStdEncoding.DecodeString(string(publicKey))
	if err != nil {
		log.Error("base64 Decode error:", err)
	}
	log.Debug("base64 key...", string(bKey))
	pub, err := x509.ParsePKCS1PublicKey(bKey)
	if err != nil {
		log.Error("ParsePKCS1PublicKey err: %s", err)
		return err
	}
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], bSignData)
}
