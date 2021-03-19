/*************************************************************************
	> File Name: redis_spider.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月10日 星期四 13时44分56秒
*************************************************************************/

package main

import (
	"github.com/sgs921107/gspider"
	"sync"
	"time"
)

// settings
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
	Timeout:        time.Second * 10,
	RedisAddr:      "172.17.0.1:6379",
	RedisDB:        6,
	RedisPassword:  "qaz123",
	RedisPrefix:    "simple",
	MaxIdleTimeout: time.Second * 10,
}

var redisKey = "start_urls"

func main() {
	var wg = &sync.WaitGroup{}
	spider := gspider.NewRedisSpider(redisKey, settings)
	wg.Add(1)
	// 向rediskey中插入url
	go func() {
		defer wg.Done()
		urls := []string{
			"http://httpbin.org/",
			"http://httpbin.org/ip",
			"http://httpbin.org/cookies/set?a=b&c=d",
			"http://httpbin.org/cookies",
		}
		for _, url := range urls {
			// 可以使用redis客户端向redisKey中添加
			spider.Client.RPush(spider.RedisKey, url)
			// 也可以直接使用spider的Queue的AddURL方法
			spider.Queue.AddURL(url)
		}
	}()
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
	wg.Wait()
}
