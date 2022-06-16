package mvutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

type AESECBEncrypt struct {
	key       []byte
	blockSize int
}

func NewAESECBEncrypt(key []byte, blockSize int) *AESECBEncrypt {
	return &AESECBEncrypt{
		key:       key,
		blockSize: blockSize,
	}
}

func (ecb *AESECBEncrypt) Encrypt(srcStr string) (val string, err error) {
	//key只能是 16 24 32长度
	block, err := aes.NewCipher([]byte(ecb.key))
	if err != nil {
		return
	}

	src := []byte(srcStr)
	//padding
	src = PKCS5Padding(src, ecb.blockSize)
	//返回加密结果
	encryptData := make([]byte, len(src))
	//分组分块加密
	encrypt := NewECBEncrypter(block)
	encrypt.CryptBlocks(encryptData, src)
	return url.QueryEscape(base64.StdEncoding.EncodeToString(encryptData)), nil
}

func (ecb *AESECBEncrypt) Decrypt(srcStr string) (val string, err error) {
	srcStr, _ = url.QueryUnescape(srcStr)
	src, err := base64.StdEncoding.DecodeString(srcStr)
	if err != nil {
		return
	}

	//key只能是 16 24 32长度
	block, err := aes.NewCipher([]byte(ecb.key))
	if err != nil {
		return
	}
	//返回解密结果
	decryptData := make([]byte, len(src))
	//分组分块加密
	decrypt := NewECBDecrypter(block)
	decrypt.CryptBlocks(decryptData, src)
	return string(PKCS5UnPadding(decryptData)), nil
}

type AESCBCEncrypt struct {
	key []byte
	iv  []byte
}

func NewAESCBCEncrypt(key, iv []byte) *AESCBCEncrypt {
	return &AESCBCEncrypt{
		key: key,
		iv:  iv,
	}
}

func (cbc *AESCBCEncrypt) Encrypt(srcStr string) (val string, err error) {
	block, err := aes.NewCipher(cbc.key)
	if err != nil {
		return
	}

	src := []byte(srcStr)
	blockSize := block.BlockSize()
	src = ZeroPadding(src, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, cbc.iv)
	dest := make([]byte, len(src))
	blockMode.CryptBlocks(dest, src)
	return hex.EncodeToString(dest), nil
}

func (cbc *AESCBCEncrypt) Decrypt(srcStr string) (val string, err error) {
	src, err := hex.DecodeString(srcStr)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(cbc.key)
	if err != nil {
		return
	}

	blockMode := cipher.NewCBCDecrypter(block, cbc.iv)
	dest := make([]byte, len(src))

	blockMode.CryptBlocks(dest, src)
	dest = ZeroUnPadding(dest)
	return string(dest), nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	for i := len(origData) - 1; ; i-- {
		if origData[i] != 0 {
			return origData[:i+1]
		}
	}
	return nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//填充
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
