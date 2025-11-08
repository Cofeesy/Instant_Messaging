package utils

import(
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5加密
func EncryptMD5(data ,salt string) string {
	hash := md5.Sum([]byte(data+salt))
	return hex.EncodeToString(hash[:])
}

// 解密,因为md5加密是单向的，所以这里的解密其实是对比是否相同
// passwd:数据库存储加密密码
// password:原有密码
func DecryptMD5(salt, passwd, password string) bool {
	encrypted := EncryptMD5(password,salt)
	return strings.EqualFold(encrypted, passwd)
}




