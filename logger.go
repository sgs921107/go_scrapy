/*************************************************************************
	> File Name: logger.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月08日 星期二 10时29分57秒
 ************************************************************************/

package gspider

import (
	"io"
	"log"
	"os"
)

type Logging struct {
	Output io.Writer
	Prefix string
	Flag   int
	logger *log.Logger
}

func (l *Logging) GetLogger() *log.Logger {
	if l.Output == nil {
		l.Output = os.Stderr
	}
	l.logger = log.New(l.Output, l.Prefix, l.Flag)
	return l.logger
}

func NewLogger(output io.Writer, prefix string, flag int) *log.Logger {
	if flag == 0 {
		flag = log.Ldate | log.Ltime
	}
	logging := &Logging{
		Output: output,
		Prefix: prefix,
		Flag:   flag,
	}
	return logging.GetLogger()
}
