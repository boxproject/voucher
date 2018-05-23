package util

import (
	"crypto/aes"
	"crypto/cipher"
	//"encoding/base64"
	"encoding/hex"
	log "github.com/alecthomas/log4go"
)

//Decrypter
func CBCDecrypter(data, cbcKey []byte) []byte {
	if len(data) < aes.BlockSize {
		log.Error("data is too short.")
		return []byte{}
	}

	deDataWithIv, err := hex.DecodeString(string(data))
	if err != nil {
		log.Error("hex Decode error:", err)
	}

	iv := make([]byte, aes.BlockSize)
	copy(iv, deDataWithIv[:aes.BlockSize])

	c, err := aes.NewCipher(cbcKey)
	if err != nil {
		log.Error("NewCipher[%s] err: %s", cbcKey, err)
	}
	decrypter := cipher.NewCBCDecrypter(c, iv)

	deData := deDataWithIv[aes.BlockSize:]

	cbcData := make([]byte, len(deData))
	decrypter.CryptBlocks(cbcData, deData)

	length := len(deData)

	pubkey := cbcData[:length-int(cbcData[length-1])]
	return pubkey

}

//Encrypter
//func CBCEncrypter(data []byte) []byte {
//	key := make([]byte, len([]byte(cbcKey)))
//	copy(key, []byte(cbcKey))
//	c, err := aes.NewCipher(key)
//	if err != nil {
//		fmt.Println("%s: NewCipher(%d bytes) = %s", err)
//	}
//	encrypter := cipher.NewCBCEncrypter(c, iv)
//
//	//规整数据长度
//	length := 0
//	if len(data)%encrypter.BlockSize() != 0 {
//		length = (len(data)/encrypter.BlockSize() + 1) * encrypter.BlockSize()
//	}
//	src := make([]byte, length)
//	copy(src, data)
//
//	deData := make([]byte, length)
//	encrypter.CryptBlocks(deData, src)
//	plaintext := base64.StdEncoding.EncodeToString(deData)
//	fmt.Printf("CBCEncrypter-------------------%s\n", plaintext)
//	return []byte(plaintext)
//}
