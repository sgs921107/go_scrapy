/*************************************************************************
	> File Name: main.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月05日 星期六 10时16分31秒
 ************************************************************************/

package gspider

import (
	"time"
)

// SpiderSettings spider settings
type SpiderSettings struct {
	Debug         	bool
	LogLevel		string 
	// 为空则输出至stdout
	LogFile       	string
	// 日志文件多久轮转一次
	RotationTime	time.Duration
	// 日志文件最大保存多久时间
	RotationMaxAge	time.Duration
	FlushOnStart  	bool // 开始前清空之前的数据
	ConcurrentReqs	int  // 并发
	MaxDepth      	int  // 最大深度
	DontFilter    	bool // 不过滤
	EnableCookies 	bool // 启用cookies
	Async         	bool // 启用异步
	KeepAlive     	bool
	Timeout       	time.Duration
	// 以下使用redis s	ider时需要配置
	RedisAddr     	string
	RedisDB       	int
	RedisPassword 	string
	RedisPrefix   	string
	MaxIdleTimeout	time.Duration // 最大闲置时间, redis spider使用 0表示一直运行
}

// DemoSettings spider setings实例的demo
var DemoSettings = SpiderSettings{
	Debug: ON,
	// 是否在启动前清空之前的数据
	FlushOnStart: OFF,
	ConcurrentReqs: 16,
	// 最大深度
	MaxDepth: 1,
	// 不过滤
	DontFilter: OFF,
	// 启用异步
	Async: OFF,
	// 不启用cookies
	EnableCookies: OFF,
	// 是否开启长连接 bool
	KeepAlive: OFF,
	// 超时
	Timeout: 30 * time.Second,
	// redis配置
	RedisAddr:     "127.0.0.1:6379",
	RedisDB:       0,
	RedisPassword: "",
	RedisPrefix:   "simple",
}
