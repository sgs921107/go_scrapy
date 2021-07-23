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
