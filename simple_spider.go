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
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

type SimpleSpider struct {
	BaseSpider
	Urls []string
	// Storage *queue.InMemoryQueueStorage
	Queue *queue.Queue
}

func (s *SimpleSpider) Start() {
	s.BaseSpider.Start()
	defer s.Close()
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
	storage := &queue.InMemoryQueueStorage{MaxSize: 10000}
	q, _ := queue.New(s.Settings.ConcurrentReqs, storage)
	s.Queue = q
	s.BaseSpider.Init()
}

func NewSimpleSpider(urls []string, settings *SpiderSettings) *SimpleSpider {
	spider := &SimpleSpider{
		BaseSpider: BaseSpider{
			Collector: colly.NewCollector(),
			Settings:  settings,
		},
		Urls: urls,
	}
	spider.Init()
	return spider
}
