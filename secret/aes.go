package secret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"time"
)

func NewKey() []byte {
	time := strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Println(time)
	h := md5.New()
	io.WriteString(h, time)
	return h.Sum(nil)
}

func AesEncrypt(encodeMsg []byte, key string) (string, error) {
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	encodeMsg = PKCS7Padding(encodeMsg, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(encodeMsg))
	blockMode.CryptBlocks(cryted, encodeMsg)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
