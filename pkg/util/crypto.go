package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

//pkcs7Padding pad
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

//pkcs7UnPadding pad reverse
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	//pad count
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

/* javascript
const CryptoJS = require('crypto-js');
const sKey = CryptoJS.enc.Utf8.parse('xBGInBHqkVvj8rVyTntz2rKvhNvrqOZj');
export function EncryptAES(s: string): string {
// key and iv use a same value
const encrypted = CryptoJS.AES.encrypt(s, sKey, {
iv: sKey,
mode: CryptoJS.mode.CBC,
padding: CryptoJS.pad.Pkcs7,
});
return encrypted.toString();
}*/

//AesEncrypt encryption use base64.StdEncoding.EncodeToString convert result []byte to string
func AesEncrypt(data []byte, keyStr string) ([]byte, error) {
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func AesEncryptStr(data string, keyStr string) (string, error) {
	encBytes, err := AesEncrypt([]byte(data), keyStr)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encBytes), nil
}

/* javascript
export function DecryptAES(s: string): string {
// key and iv use a same value
const decrypted = CryptoJS.AES.decrypt(s, sKey, {
iv: sKey,
mode: CryptoJS.mode.CBC,
padding: CryptoJS.pad.Pkcs7,
});
return decrypted.toString(CryptoJS.enc.Utf8);
}*/

//AesDecrypt decryption use base64.StdEncoding.DecodeString convert data string to []byte
func AesDecrypt(data []byte, keyStr string) ([]byte, error) {
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

func AesDecryptStr(data string, keyStr string) (string, error) {
	decBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	decBytes, err = AesDecrypt(decBytes, keyStr)
	if err != nil {
		return "", err
	}
	return string(decBytes), nil
}

// DesEncrypt DES CBC PKCS7 encryption
func DesEncrypt(dataStr, keyStr string) (string, error) {
	src := []byte(dataStr)
	key := []byte((keyStr + "--------")[0:8])
	iv := key

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()

	// pkcsPadding pkcs5/pkcs7补码算法
	pkcs7Padding(src, bs)

	blocMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(src))
	blocMode.CryptBlocks(encrypted, src)

	return strings.ToUpper(hex.EncodeToString(encrypted)), nil
}

// MD5Encrypt md5 encryption
func MD5Encrypt(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

// Sha256 encryption
func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	encodedStr := hex.EncodeToString(bs)
	return encodedStr
}
