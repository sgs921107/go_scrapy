/*************************************************************************
	> File Name: main.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月05日 星期六 10时16分31秒
 ************************************************************************/

package gspiders

type RedisParamsType struct {
}

type SpiderSettings struct {
	Debug          bool
	LogFile        string
	LogPrefix      string
	LogFlag        int
	FlushOnStart   bool // 开始前清空之前的数据
	UserAgent      string
	ConcurrentReqs int  // 并发
	MaxDepth       int  // 最大深度
	DontFilter     bool // 不过滤
	EnableCookies  bool // 启用cookies
	Async          bool // 启用异步
	KeepAlive      bool
	Timeout        int
	MaxConns       int
	// 以下使用redis spider时需要配置
	RedisAddr      string
	RedisDB        int
	RedisPassword  string
	RedisPrefix    string
	MaxIdleTimeout int // 最大闲置时间, redis spider使用 0表示一直运行
}

// spider setings实例的demo
var DemoSettings = SpiderSettings{
	Debug: ON,
	// 是否在启动前清空之前的数据
	FlushOnStart: OFF,
	// UserAgent bool
	UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
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
	// 超时  单位：秒
	Timeout: 30,
	// 最大连接数
	MaxConns: 100,
	// redis配置
	RedisAddr:     "127.0.0.1:6379",
	RedisDB:       0,
	RedisPassword: "",
	RedisPrefix:   "simple",
}
