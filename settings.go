/*************************************************************************
	> File Name: main.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月05日 星期六 10时16分31秒
 ************************************************************************/

package gspider

// SpiderSettings spider settings
type SpiderSettings struct {
	Spider struct {
		Debug         	bool 	`default:"true"`  // 是否开启debug模式
		FlushOnStart  	bool 	`default:"true"`  // 开始前清空之前的数据
		ConcurrentReqs	int  	`default:"16"`  // 并发
		MaxDepth      	int  	`default:"10"`  // 最大深度
		DontFilter    	bool 	`default:"true"`  // 不过滤
		EnableCookies 	bool 	`default:"false"`  // 启用cookies
		Async         	bool 	`default:"true"`  // 启用异步
		KeepAlive     	bool 	`default:"false"`  // 开启长连接
		Timeout       	int64	`default:"30"`  // 超时 单位: 秒
		MaxIdleTimeout	int64 	`default:"60"`	// 最大闲置时间, redis spider使用 0表示一直运行 单位: 秒
	}
	Log	struct {
		// 日志等级
		Level string 			`default:"DEBUG"`
		// 日志文件 为空则输出至stdout
		File	string			`default:""`
		// 日志文件多久轮转一次 单位: hour
		RotationTime	uint	`default:"24"`
		// 日志文件最大保存多久时间	单位: hour
		RotationMaxAge	uint	`default:"168"`
	}
	// 以下使用redis spider时需要配置
	Redis struct {
		Prefix   	string 	`default:"gs"`
		Addr     	string 	`default:"127.0.0.1:6379"`
		DB       	int    	`default:"0"`
		Password 	string 	`default:""`
	}
}
