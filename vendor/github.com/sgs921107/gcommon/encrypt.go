/*************************************************************************
	> File Name: encrypt.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2021年06月11日 星期五 15时53分14秒
*************************************************************************/

package gcommon

import (
	"crypto/md5"
	"encoding/hex"
)

// EncryptMD5 encrypt text with md5
func EncryptMD5(text string) string {
	m := md5.New()
	m.Write([]byte(text))
	return hex.EncodeToString(m.Sum(nil))
}
