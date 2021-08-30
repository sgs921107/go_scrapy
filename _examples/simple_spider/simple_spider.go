/*************************************************************************
	> File Name: redis_spider.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月10日 星期四 13时44分56秒
*************************************************************************/

package main

import (
	"github.com/sgs921107/gspider"
	"github.com/sgs921107/gcommon"
	"sync"
)

var (
	settings = gspider.SpiderSettings{}
	settingsOnce sync.Once
)


func newSettings() gspider.SpiderSettings {
	settingsOnce.Do(func(){
		gcommon.LoadEnvFile("env_demo", true)
		gcommon.EnvIgnorePrefix()
		gcommon.EnvFill(&settings)
	})
	return settings
}

func main() {
	settings := newSettings()
	urls := []string{
		"http://httpbin.org/",
		"http://httpbin.org/ip",
		"http://httpbin.org/cookies/set?a=b&c=d",
		"http://httpbin.org/cookies",
	}
	spider := gspider.NewSimpleSpider(urls, settings)
	spider.OnRequest(func(r *gspider.Request) {
		spider.Logger.Infow("send a req",
			"method", r.Method,
			"url", r.URL,
		)
	})
	spider.OnResponse(func(r *gspider.Response) {
		spider.Logger.Infow("recv a resp",
			"url", r.Request.URL,
		)
	})
	spider.Start()
}
