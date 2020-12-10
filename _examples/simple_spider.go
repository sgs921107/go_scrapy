/*************************************************************************
	> File Name: redis_spider.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月10日 星期四 13时44分56秒
*************************************************************************/

package main

import (
	"gspider"
)

var settings = &gspider.SpiderSettings{
	Debug: true,
	// 是否在启动前清空之前的数据
	FlushOnStart: true,
	// UserAgent bool
	UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
	ConcurrentReqs: 16,
	// 最大深度
	MaxDepth: 1,
	// 不过滤
	DontFilter: false,
	// 启用异步
	Async: true,
	// 不启用cookies
	EnableCookies: false,
	// 是否开启长连接 bool
	KeepAlive: false,
	// 超时  单位：秒
	Timeout: 30,
	// 最大连接数
	MaxConns: 100,
}

var redisKey = "start_urls"

func main() {
	urls := []string{
		"http://httpbin.org/",
		"http://httpbin.org/ip",
		"http://httpbin.org/cookies/set?a=b&c=d",
		"http://httpbin.org/cookies",
	}
	spider := gspider.NewSimpleSpider(urls, settings)
	spider.OnRequest(func(r *gspider.Request) {
		spider.Logger.Printf("create a task: %s %s", r.Method, r.URL)
	})
	spider.OnResponse(func(r *gspider.Response) {
		spider.Logger.Printf("recv a resp: %s", r.Request.URL)
	})
	spider.Start()
}
