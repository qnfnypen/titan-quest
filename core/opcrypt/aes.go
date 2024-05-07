package opcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// ECB: 电子密码文本模式，最基本的工作模式，将待处理信息分组，每组分别进行加密或解密处理
// CBC: 密码分组链接模式，每个明文块先与其前一个密文块进行异或，然后再进行加密
// CFB: 密文反馈模式，前一个密文使用密钥key再加密后，与明文异或，得到密文。第一个密文需要初始向量IV加密得到。
// 		解密也同样使用加密器进行解密

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	// 判断缺少几位长度，最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	// 补足位数，复制padding
	padText := bytes.Repeat([]byte{(byte(padding))}, padding)

	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误")
	}
	// 获取填充的个数
	unPadding := int(data[length-1])

	return data[:length-unPadding], nil
}

// AesEncryptCBC 加密后再进行base64编码
func AesEncryptCBC(data []byte, key []byte) (string, error) {
	// 创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// 判断加密块的大小
	blockSize := block.BlockSize()
	// 填充
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

// AesDecryptCBC 解密
func AesDecryptCBC(data string, key []byte) (crypted []byte, err error) {
	// 处理 panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("aes decrypt error:%v", r)
			return
		}
	}()

	// base64解码
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// 创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 判断加密块的大小
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted = make([]byte, len(dataByte))
	blockMode.CryptBlocks(crypted, dataByte)
	// 去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}

	return crypted, nil
}