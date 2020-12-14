/*************************************************************************
	> File Name: simple.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月07日 星期一 10时23分01秒
 ************************************************************************/
/*
simple spider
基于colly实现类似python中scrapy中的spider
*/

package gspider

import (
	"github.com/gocolly/colly/v2"
)

type SimpleSpider struct {
	BaseSpider
	Urls  []string
	Queue *Queue
}

func (s *SimpleSpider) Start() {
	s.BaseSpider.Start()
	defer func() {
		s.Queue.Stop()
		s.Close()
	}()
	for _, url := range s.Urls {
		s.Queue.AddURL(url)
	}
	s.Queue.Run(s.Collector)
	s.Collector.Wait()
}

func (s *SimpleSpider) Close() {
	s.BaseSpider.Close()
}

func (s *SimpleSpider) Init() {
	storage := &InMemoryQueueStorage{MaxSize: 10000}
	q, _ := NewQueue(s.settings.ConcurrentReqs, storage)
	s.Queue = q
	s.BaseSpider.Init()
}

func NewSimpleSpider(urls []string, settings *SpiderSettings) *SimpleSpider {
	spider := &SimpleSpider{
		BaseSpider: BaseSpider{
			Collector: colly.NewCollector(),
			settings:  settings,
		},
		Urls: urls,
	}
	spider.Init()
	return spider
}
