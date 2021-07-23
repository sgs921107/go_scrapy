/*************************************************************************
	> File Name: Redis.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月07日 星期一 10时23分01秒
 ************************************************************************/
/*
redis spider
类似python的scrapy-redis中的redis_spider
启动后会监听start url队列中的任务进行下载
*/

package gspider

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/redisstorage"
	// "runtime"
	"github.com/sgs921107/gcommon"
	"github.com/sgs921107/gredis"
	"sync"
	"sync/atomic"
	"time"
)

/*
RedisSpider redis spider
*/
type RedisSpider struct {
	BaseSpider
	RedisKey string
	Client   *gredis.Client
	Queue    *Queue
	last     int64
	wg       *sync.WaitGroup
}

/*
listenStartURLs监听start_urls队列
*/
func (s *RedisSpider) listenStartURLs() {
	defer s.wg.Done()
	for {
		if atomic.LoadUint32(&s.exit) != 0 {
			break
		}
		if url, err := s.Client.LPop(s.RedisKey).Result(); err == nil {
			s.Queue.AddURL(url)
		} else {
			time.Sleep(500 * time.Millisecond)
			// runtime.Gosched()
		}
	}
}

/*
Start 启动redis spider
*/
func (s *RedisSpider) Start() {
	s.BaseSpider.Start()
	defer s.Close()
	s.wg.Add(1)
	go s.listenStartURLs()
	for {
		s.Queue.Run(s.Collector)
		s.Collector.Wait()
		if s.settings.Spider.MaxIdleTimeout != 0 {
			// 纳秒时间戳
			now := gcommon.TimeStamp(3)
			// 超出最大闲置时间则退出
			if now-atomic.LoadInt64(&s.last) > int64(s.settings.Spider.MaxIdleTimeout) * int64(time.Second) {
				break
			}
		}
		// 重置queue的状态,等待下一次启动
		s.Queue.Stop()
		time.Sleep(500 * time.Millisecond)
		// runtime.Gosched()
	}
}

/*
recordLastTime 记录最后一个请求发出的时间
*/
func (s *RedisSpider) recordLastTime(*Request) {
	atomic.StoreInt64(&s.last, gcommon.TimeStamp(3))
}

/*
Close 释放资源
*/
func (s *RedisSpider) Close() {
	s.Client.Close()
	s.BaseSpider.Close()
	// 等待监听start urls队列的任务结束
	s.wg.Wait()
}

/*
Init 配置使用redis存储
*/
func (s *RedisSpider) Init() {
	s.RedisKey = s.settings.Redis.Prefix + ":" + "start_urls"
	storage := &redisstorage.Storage{
		Address:  s.settings.Redis.Addr,
		Password: s.settings.Redis.Password,
		DB:       s.settings.Redis.DB,
		Prefix:   s.settings.Redis.Prefix,
	}
	err := s.Collector.SetStorage(storage)
	// 下面使用到logger 需先init base spider
	// 不能在set storage前执行，会导致disable cookies被覆盖
	s.BaseSpider.Init()
	if err != nil {
		s.Logger.WithFields(LogFields{
			"errMsg": err.Error(),
		}).Fatal("Set redis storage failed")
	}
	s.Client = gredis.NewClientFromRedisClient(storage.Client)
	if s.settings.Spider.FlushOnStart {
		if err := storage.Clear(); err != nil {
			s.Logger.WithFields(LogFields{
				"errMsg": err.Error(),
			}).Error("clear previous data of redis storage failed")
		}
		s.Client.Del(s.RedisKey)
	}
	q, _ := NewQueue(s.settings.Spider.ConcurrentReqs, storage)
	s.Queue = q
	// 如果配置了最大闲置时间
	if s.settings.Spider.MaxIdleTimeout != 0 {
		s.OnRequest(s.recordLastTime)
		atomic.StoreInt64(&s.last, gcommon.TimeStamp(3))
	}
}

/*
NewRedisSpider 生成一个redis spider实例
*/
func NewRedisSpider(settings SpiderSettings) *RedisSpider {
	spider := &RedisSpider{
		BaseSpider: BaseSpider{
			Collector: colly.NewCollector(),
			settings:  settings,
		},
		wg:       &sync.WaitGroup{},
	}
	spider.Init()
	return spider
}
