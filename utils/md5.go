package utils

import(
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 加密
func EncryptMD5(data ,salt string) string {
	hash := md5.Sum([]byte(data+salt))
	return hex.EncodeToString(hash[:])
}

// 解密,因为md5加密是单向的，所以这里的解密其实是对比是否相同
// passwd:待验证密码
// password:原有密码
func DecryptMD5(salt, passwd, password string) bool {
	encrypted := EncryptMD5(password,salt)
	// println(encrypted)
	// println(ConcatMD5(encrypted, salt))
	// println(passwd)
	return strings.EqualFold(encrypted, passwd)
}

// 字符串拼接
// func ConcatMD5(encrypted, salt string) string {
// 	return strings.Join([]string{encrypted, salt}, "")
// }


// utils.DecryptMD5(user.Salt, user.Password, utils.ConcatMD5(utils.EncryptMD5(passwd), user.Salt))
// func main() {
// 	// 测试
// 	salt := "123456"
// 	password := "password"
// 	encrypted := encryptMD5(password)
// 	println("原始字符串:", salt)
// 	println("加密后字符串:", encrypted)
// 	println("拼接后字符串:", concatMD5(encrypted, salt))
// 	println("解密结果:", decryptMD5(salt, concatMD5(encrypted, salt), password))
// }



