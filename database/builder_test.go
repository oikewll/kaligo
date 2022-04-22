package database

import (
    // "encoding/json"
    "os"
    "log"
    "testing"
    // "github.com/stretchr/testify/assert"
)

// 每个测试函数执行之前会先执行这个，用于初始化，每个测试都是纯净的
func TestMain(m *testing.M) {
    // config.Set("cache.config.driver", "memory")
    log.Printf("%v", "TestMain")
    os.Exit(m.Run())
}

// func TestNew(t *testing.T) {
//     assert.NotNil(t, mux)
// }
//
// // 测试padding  unpadding
// func TestPadding(t *testing.T) {
//     tool := NewAesTool([]byte{}, 16)
//     src := []byte{1, 2, 3, 4, 5}
//     src = tool.padding(src)
//     // log.Printf("%v", src)
//     src = tool.unPadding(src)
//     // log.Printf("%v", src)
// }
//
// // 测试 AES ECB 加密解密
// func TestEncryptDecrypt(t *testing.T) {
//     key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}
//     blickSize := 16
//     tool := NewAesTool(key, blickSize)
//     encryptStr := "32334erew3232334erew3232334erew3232334erew3232334erew3232334erew3232334erew3232334erew323232334erew3232334erew32"
//     // log.Printf("%v", len(encryptStr))
//     encryptData, _ := tool.Encrypt([]byte(encryptStr))
//     // log.Printf("%v", len(encryptData))
//     decryptData, _ := tool.Decrypt(encryptData)
//     log.Printf("%v", len(decryptData))
// }

