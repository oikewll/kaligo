package util

import (
    "bytes"
    "crypto/aes"
)

// AES ECB模式的加密解密
type AesTool struct {
    // 128 192  256位的其中一个 长度 对应分别是 16 24 32 字节长度
    Key       []byte
    BlockSize int
}

func NewAesTool(key []byte, blockSize int) *AesTool {
    return &AesTool{Key: key, BlockSize: blockSize}
}

func (this *AesTool) padding(src []byte) []byte {
    // 填充个数
    paddingCount := aes.BlockSize - len(src)%aes.BlockSize
    if paddingCount == 0 {
        return src
    } else {
        // 填充数据
        return append(src, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
    }
}

func (this *AesTool) unPadding(src []byte) []byte {
    for i := len(src) - 1; ; i-- {
        if src[i] != 0 {
            return src[:i+1]
        } 
    }
}

func (this *AesTool) Encrypt(src []byte) ([]byte, error) {
    // key 只能是 16 24 32 长度
    block, err := aes.NewCipher([]byte(this.Key))
    if err != nil {
        return nil, err
    }
    src = this.padding(src)
    // 返回加密结果
    encryptData := make([]byte, len(src))
    // 存储每次加密的数据
    tmpData := make([]byte, this.BlockSize)

    // 分组分块加密
    for index := 0; index < len(src); index += this.BlockSize {
        block.Encrypt(tmpData, src[index:index+this.BlockSize])
        copy(encryptData, tmpData)
    }
    return encryptData, nil
}

func (this *AesTool) Decrypt(src []byte) ([]byte, error) {
    // key 只能是 16 24 32 长度
    block, err := aes.NewCipher([]byte(this.Key))
    if err != nil {
        return nil, err
    }
    // 返回加密结果
    decryptData := make([]byte, len(src))
    // 存储每次加密的数据
    tmpData := make([]byte, this.BlockSize)

    // 分组分块加密
    for index := 0; index < len(src); index += this.BlockSize {
        block.Decrypt(tmpData, src[index:index+this.BlockSize])
        copy(decryptData, tmpData)
    }
    return this.unPadding(decryptData), nil
}
