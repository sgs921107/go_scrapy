/*************************************************************************
	> File Name: funcs.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 19时39分01秒
 ************************************************************************/

package gcommon

import (
	"net/url"
	"os"
	"time"
)

// TimeStampFlag 定义类型为int
type TimeStampFlag int

const (
	// SECOND 0: 秒
	SECOND TimeStampFlag = iota
	// MILLISECOND 1: 毫秒
	MILLISECOND
	// MICROSECOND 2: 微秒
	MICROSECOND
	// NANOSECOND 3: 纳秒
	NANOSECOND
)

/*
TimeStamp 获取时间戳
接收一个整形 0-3
0-秒, 1-毫秒, 2-微妙, 3-纳秒
*/
func TimeStamp(flag int) int64 {
	now := time.Now()
	switch TimeStampFlag(flag) {
	case SECOND:
		return now.Unix()
	case MILLISECOND:
		return now.UnixNano() / 1e6
	case MICROSECOND:
		return now.UnixNano() / 1e3
	case NANOSECOND:
		return now.UnixNano()
	default:
		return 0
	}
}

/*
FetchURLHost 提取url的host
*/
func FetchURLHost(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

/*
PathIsExist 判断路径是否存在
*/
func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
