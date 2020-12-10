/*************************************************************************
	> File Name: spiders.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月07日 星期一 10时22分35秒
 ************************************************************************/

package gspider

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

// 起别名
type (
	Request          = colly.Request
	Response         = colly.Response
	Context          = colly.Context
	ProxyFunc        = colly.ProxyFunc
	RequestCallback  = colly.RequestCallback
	ResponseCallback = colly.ResponseCallback
	HTMLCallback     = colly.HTMLCallback
	XMLCallback      = colly.XMLCallback
	ScrapedCallback  = colly.ScrapedCallback
)

// spider结构
type BaseSpider struct {
	Collector      *colly.Collector
	Settings       *SpiderSettings
	Logger         *log.Logger
	output         *os.File
	curReqCounter  int64
	curRespCounter int64
	reqCounter     int64
	respCounter    int64
	exit           uint32
}

func (s *BaseSpider) Start() {
	s.Logger.Print("==========================spider start====================================")
}

// 给请求计数器追加1
func (s *BaseSpider) recordReq(*Request) {
	atomic.AddInt64(&s.curReqCounter, 1)
	atomic.AddInt64(&s.reqCounter, 1)
}

// 给请求计数器追加1
func (s *BaseSpider) recordResp(*Response) {
	atomic.AddInt64(&s.curRespCounter, 1)
	atomic.AddInt64(&s.respCounter, 1)
}

func (s *BaseSpider) showCounter() {
	ticker := time.NewTicker(time.Minute)
	for {
		if atomic.LoadUint32(&s.exit) != 0 {
			break
		}
		select {
		case <-ticker.C:
			s.Logger.Printf(
				"----------------------crawled (%d/%d), (%d/%d)/min------------------------",
				atomic.LoadInt64(&s.reqCounter),
				atomic.LoadInt64(&s.respCounter),
				atomic.SwapInt64(&s.curReqCounter, 0),
				atomic.SwapInt64(&s.curRespCounter, 0),
			)
		}
	}
}

// 配置http配置
func (s *BaseSpider) SetHttp() {
	s.Collector.WithTransport(&http.Transport{
		DisableKeepAlives: s.Settings.KeepAlive,
		MaxIdleConns:      s.Settings.MaxConns,
	})
}

// 配置扩展 自动添加user agent、referer
func (s *BaseSpider) SetExtensions() {
	// 自动添加随机生成的user agent
	extensions.RandomUserAgent(s.Collector)
	// 添加referer信息
	extensions.Referer(s.Collector)
}

func (s *BaseSpider) OnRequest(f RequestCallback) {
	s.Collector.OnRequest(f)
}

func (s *BaseSpider) OnError(f func(r *Response, err error)) {
	s.Collector.OnError(f)
}

// func (s *BaseSpider) OnResponseHeaders(f colly.ResponseHeadersCallback) {
// 	s.Collector.OnResponseHeaders(f)
// }

func (s *BaseSpider) OnResponse(f ResponseCallback) {
	s.Collector.OnResponse(f)
}

func (s *BaseSpider) OnHTML(goquerySelector string, f HTMLCallback) {
	s.Collector.OnHTML(goquerySelector, f)
}

func (s *BaseSpider) OnXML(goquerySelector string, f XMLCallback) {
	s.Collector.OnXML(goquerySelector, f)
}

func (s *BaseSpider) OnScraped(f ScrapedCallback) {
	s.Collector.OnScraped(f)
}

func (s *BaseSpider) Post(url string, data map[string]string) error {
	if err := s.Collector.Post(url, data); err != nil {
		s.Logger.Printf("HttpError: url: %s, data %v, err msg: %s", url, data, err.Error())
		return err
	}
	return nil
}

func (s *BaseSpider) PostMult(url string, data map[string][]byte) error {
	if err := s.Collector.PostMultipart(url, data); err != nil {
		s.Logger.Printf("HttpError: url: %s, data %v, err msg: %s", url, data, err.Error())
		return err
	}
	return nil
}

func (s *BaseSpider) Cookies(url string) []*http.Cookie {
	return s.Collector.Cookies(url)
}

func (s *BaseSpider) SetProxy(proxyURL string) error {
	return s.Collector.SetProxy(proxyURL)
}

func (s *BaseSpider) SetProxyFunc(f ProxyFunc) {
	s.Collector.SetProxyFunc(f)
}

// 给spider配置一个logger
func (s *BaseSpider) SetLogger() {
	if s.Settings.LogFile != "" {
		output, err := os.OpenFile(s.Settings.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		s.output = output
	} else {
		s.output = os.Stderr
	}
	prefix := s.Settings.LogPrefix
	flag := s.Settings.LogFlag
	s.Logger = NewLogger(s.output, prefix, flag)
	// 配置debugger
	if s.Settings.Debug == true {
		s.Collector.SetDebugger(&debug.LogDebugger{
			Output: s.output,
			Prefix: prefix,
			Flag:   flag,
		})
	}
}

// 加载配置
func (s *BaseSpider) LoadSettings() {
	// 配置最大深度
	s.Collector.MaxDepth = s.Settings.MaxDepth
	// 配置是否可重复抓取
	s.Collector.AllowURLRevisit = s.Settings.DontFilter
	// http 配置
	s.SetHttp()
	// 配置是否启用异步
	s.Collector.Async = s.Settings.Async
	// 设置timeout
	s.Collector.SetRequestTimeout(time.Duration(s.Settings.Timeout) * time.Second)
	// 配置是否启用cookies
	if s.Settings.EnableCookies == OFF {
		s.Collector.DisableCookies()
	}
	s.SetLogger()
}

// 初始化工作
func (s *BaseSpider) Init() {
	s.LoadSettings()
	s.SetExtensions()
	s.OnRequest(s.recordReq)
	s.OnResponse(s.recordResp)
	go s.showCounter()
	// s.OnError(func(r *colly.Response, err error) {
	// 	log.Printf("HttpError: url: %s, code: %d, err msg: %s", r.Request.URL, r.StatusCode, err.Error())
	// })
}

// 释放资源
func (s *BaseSpider) Close() {
	atomic.StoreUint32(&s.exit, 1)
	s.Logger.Print("==========================spider close====================================")
	if s.output != nil {
		s.output.Close()
	}
}

func init() {
	// 配置go使用的cpu数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
