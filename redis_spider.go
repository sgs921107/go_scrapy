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
		if s.settings.MaxIdleTimeout != 0 {
			now := gcommon.TimeStamp(0)
			// 超出最大闲置时间则退出
			maxIdleTimeout := int64(s.settings.MaxIdleTimeout)
			// 如果最大闲置时间配置过小，保证所有发出的请求已结束
			if maxIdleTimeout <= int64(s.settings.Timeout) {
				maxIdleTimeout += int64(s.settings.Timeout)
			}
			if now-atomic.LoadInt64(&s.last) > maxIdleTimeout {
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
	atomic.StoreInt64(&s.last, gcommon.TimeStamp(0))
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
	storage := &redisstorage.Storage{
		Address:  s.settings.RedisAddr,
		Password: s.settings.RedisPassword,
		DB:       s.settings.RedisDB,
		Prefix:   s.settings.RedisPrefix,
	}
	err := s.Collector.SetStorage(storage)
	if err != nil {
		s.Logger.Fatalf("set redis storage failed: %s", err.Error())
		panic(err)
	}
	s.Client = gredis.NewClientFromRedisClient(storage.Client)
	if s.settings.FlushOnStart {
		if err := storage.Clear(); err != nil {
			s.Logger.Fatal("clear previous data of redis storage failed: " + err.Error())
		}
	}
	q, _ := NewQueue(s.settings.ConcurrentReqs, storage)
	s.Queue = q
	// 如果配置了最大闲置时间
	if s.settings.MaxIdleTimeout != 0 {
		s.OnRequest(s.recordLastTime)
		atomic.StoreInt64(&s.last, gcommon.TimeStamp(0))
	}
	s.BaseSpider.Init()
}

/*
NewRedisSpider 生成一个redis spider实例
*/
func NewRedisSpider(redisKey string, settings *SpiderSettings) *RedisSpider {
	// default redisKey
	if redisKey == "" {
		redisKey = "start_urls"
	}
	// reids key的prefix
	prefix := settings.RedisPrefix
	// 给redisKey 添加前缀
	redisKey = prefix + ":" + redisKey
	spider := &RedisSpider{
		BaseSpider: BaseSpider{
			Collector: colly.NewCollector(),
			settings:  settings,
		},
		RedisKey: redisKey,
		wg:       &sync.WaitGroup{},
	}
	spider.Init()
	return spider
}
