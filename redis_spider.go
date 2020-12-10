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
	"github.com/go-redis/redis"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	// "runtime"
	"github.com/sgs921107/gcommon"
	"sync/atomic"
	"time"
)

// redis spider
type RedisSpider struct {
	BaseSpider
	RedisKey string
	// Storage  *redisstorage.Storage
	Client *redis.Client
	Queue  *queue.Queue
	last   int64
}

// 监听start_urls队列
func (s *RedisSpider) ListenStartUrls() {
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

// 启动redis spider
func (s *RedisSpider) Start() {
	s.BaseSpider.Start()
	defer s.Close()
	go s.ListenStartUrls()
	for {
		s.Queue.Run(s.Collector)
		s.Collector.Wait()
		if s.Settings.MaxIdleTimeout != 0 {
			now := gcommon.TimeStamp(0)
			// 超出最大闲置时间则退出
			maxIdleTimeout := int64(s.Settings.MaxIdleTimeout)
			// 如果最大闲置时间配置过小，保证所有发出的请求已结束
			if maxIdleTimeout <= int64(s.Settings.Timeout) {
				maxIdleTimeout += int64(s.Settings.Timeout)
			}
			if now-atomic.LoadInt64(&s.last) > maxIdleTimeout {
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
		// runtime.Gosched()
	}
}

func (s *RedisSpider) recordLastTime(*Request) {
	atomic.StoreInt64(&s.last, gcommon.TimeStamp(0))
}

// 释放资源
func (s *RedisSpider) Close() {
	s.Client.Close()
	s.BaseSpider.Close()
}

// 配置使用redis存储
func (s *RedisSpider) Init() {
	storage := &redisstorage.Storage{
		Address:  s.Settings.RedisAddr,
		Password: s.Settings.RedisPassword,
		DB:       s.Settings.RedisDB,
		Prefix:   s.Settings.RedisPrefix,
	}
	err := s.Collector.SetStorage(storage)
	if err != nil {
		s.Logger.Fatalf("set redis storage failed: %s", err.Error())
		panic(err)
	}
	s.Client = storage.Client
	if s.Settings.FlushOnStart {
		if err := storage.Clear(); err != nil {
			s.Logger.Fatal("clear previous data of redis storage failed: " + err.Error())
		}
	}
	q, _ := queue.New(s.Settings.ConcurrentReqs, storage)
	s.Queue = q
	// 如果配置了最大闲置时间
	if s.Settings.MaxIdleTimeout != 0 {
		s.OnRequest(s.recordLastTime)
		atomic.StoreInt64(&s.last, gcommon.TimeStamp(0))
	}
	s.BaseSpider.Init()
}

// 生成一个redis spider实例
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
			Settings:  settings,
		},
		RedisKey: redisKey,
	}
	spider.Init()
	return spider
}
