/*************************************************************************
	> File Name: redis_spider.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月10日 星期四 13时44分56秒
*************************************************************************/

package main

import (
	"github.com/sgs921107/gspider"
	"time"
)

var settings = &gspider.SpiderSettings{
	Debug: true,
	// 是否在启动前清空之前的数据
	FlushOnStart:   true,
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
	Timeout: 30 * time.Second,
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
		spider.Logger.WithFields(gspider.LogFields{
			"method": r.Method,
			"url":    r.URL,
		}).Info("send a req")
	})
	spider.OnResponse(func(r *gspider.Response) {
		spider.Logger.WithFields(gspider.LogFields{
			"url": r.Request.URL,
		}).Info("recv a resp")
	})
	spider.Start()
}
