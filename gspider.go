/*************************************************************************
	> File Name: spiders.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月07日 星期一 10时22分35秒
 ************************************************************************/

package gspider

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/sgs921107/glogging"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

// 起别名
type (
	// Request          = colly.Request
	Request          = colly.Request
	// Response         = colly.Response
	Response         = colly.Response
	// Context          = colly.Context
	Context          = colly.Context
	// ProxyFunc        = colly.ProxyFunc
	ProxyFunc        = colly.ProxyFunc
	// RequestCallback  = colly.RequestCallback
	RequestCallback  = colly.RequestCallback
	// ResponseCallback = colly.ResponseCallback
	ResponseCallback = colly.ResponseCallback
	// HTMLCallback     = colly.HTMLCallback
	HTMLCallback     = colly.HTMLCallback
	// XMLCallback      = colly.XMLCallback
	XMLCallback      = colly.XMLCallback
	// ScrapedCallback  = colly.ScrapedCallback
	ScrapedCallback  = colly.ScrapedCallback
	// LogFields logrus fields
	LogFields		 = glogging.Fields
	// Logger logrus logger
	Logger			 = glogging.Logger
)

// BaseSpider spider结构
type BaseSpider struct {
	Collector      *colly.Collector
	settings       *SpiderSettings
	Logger         *Logger
	curReqCounter  int64
	curRespCounter int64
	reqCounter     int64
	respCounter    int64
	exit           uint32
}

// Start 启动方法
func (s *BaseSpider) Start() {
	s.Logger.Info("==========================spider start====================================")
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
			s.Logger.Infof(
				"----------------------crawled (%d/%d), (%d/%d)/min------------------------",
				atomic.LoadInt64(&s.reqCounter),
				atomic.LoadInt64(&s.respCounter),
				atomic.SwapInt64(&s.curReqCounter, 0),
				atomic.SwapInt64(&s.curRespCounter, 0),
			)
		}
	}
}

// SetHTTP http配置
func (s *BaseSpider) SetHTTP(transport *http.Transport) {
	s.Collector.WithTransport(transport)
}

// SetExtensions 配置扩展 自动添加user agent、referer
func (s *BaseSpider) SetExtensions() {
	// 自动添加随机生成的user agent
	extensions.RandomUserAgent(s.Collector)
	// 添加referer信息
	extensions.Referer(s.Collector)
}

// OnRequest 发送请求前执行
func (s *BaseSpider) OnRequest(f RequestCallback) {
	s.Collector.OnRequest(f)
}

// OnError 出错时执行
func (s *BaseSpider) OnError(f func(r *Response, err error)) {
	s.Collector.OnError(f)
}

// OnResponseHeaders 接收resp headers后执行
func (s *BaseSpider) OnResponseHeaders(f colly.ResponseHeadersCallback) {
	s.Collector.OnResponseHeaders(f)
}

// OnResponse 接收resp后执行
func (s *BaseSpider) OnResponse(f ResponseCallback) {
	s.Collector.OnResponse(f)
}

// OnHTML resp是html
func (s *BaseSpider) OnHTML(goquerySelector string, f HTMLCallback) {
	s.Collector.OnHTML(goquerySelector, f)
}

// OnXML 页面是XML时执行
func (s *BaseSpider) OnXML(goquerySelector string, f XMLCallback) {
	s.Collector.OnXML(goquerySelector, f)
}

// OnScraped 请求任务结束后执行
func (s *BaseSpider) OnScraped(f ScrapedCallback) {
	s.Collector.OnScraped(f)
}

// Post 发送一个post请求
func (s *BaseSpider) Post(url string, data map[string]string) error {
	return s.Collector.Post(url, data)
}

// PostMult 发送一个post请求
func (s *BaseSpider) PostMult(url string, data map[string][]byte) error {
	return s.Collector.PostMultipart(url, data)
}

// Cookies 获取cookies
func (s *BaseSpider) Cookies(url string) []*http.Cookie {
	return s.Collector.Cookies(url)
}

// SetProxy 设置代理
func (s *BaseSpider) SetProxy(proxyURL string) error {
	return s.Collector.SetProxy(proxyURL)
}

// SetProxyFunc 设置代理方法
func (s *BaseSpider) SetProxyFunc(f ProxyFunc) {
	s.Collector.SetProxyFunc(f)
}

// SetLogger 给spider配置一个logger
func (s *BaseSpider) SetLogger() {
	s.Logger = glogging.NewLogging(&glogging.Options{
		Level: s.settings.LogLevel,
		FilePath: s.settings.LogFile,
		RotationTime: s.settings.RotationTime,
		RotationMaxAge: s.settings.RotationMaxAge,
	}).GetLogger()
	// 配置debugger
	if s.settings.Debug == true {
		s.Collector.SetDebugger(&debug.LogDebugger{
			Output: s.Logger.Out,
		})
	}
}

// LoadSettings 加载配置
func (s *BaseSpider) LoadSettings() {
	// 配置最大深度
	s.Collector.MaxDepth = s.settings.MaxDepth
	// 配置是否可重复抓取
	s.Collector.AllowURLRevisit = s.settings.DontFilter
	transport := &http.Transport{
		DisableKeepAlives: s.settings.KeepAlive,
		MaxIdleConns:      s.settings.MaxConns,
	}
	// http 配置
	s.SetHTTP(transport)
	// 配置是否启用异步
	s.Collector.Async = s.settings.Async
	// 设置timeout
	s.Collector.SetRequestTimeout(s.settings.Timeout)
	// 配置是否启用cookies
	if s.settings.EnableCookies == OFF {
		s.Collector.DisableCookies()
	}
	s.SetLogger()
}

// Init 初始化工作
func (s *BaseSpider) Init() {
	// 如果最大闲置时间配置过小，保证所有发出的请求已结束
	if s.settings.MaxIdleTimeout != 0 && s.settings.MaxIdleTimeout <= s.settings.Timeout {
		s.settings.MaxIdleTimeout += s.settings.Timeout
	}
	s.LoadSettings()
	s.SetExtensions()
	s.OnRequest(s.recordReq)
	s.OnResponse(s.recordResp)
	go s.showCounter()
	// s.OnError(func(r *colly.Response, err error) {
	// 	log.Printf("HttpError: url: %s, code: %d, err msg: %s", r.Request.URL, r.StatusCode, err.Error())
	// })
}

// Close 释放资源
func (s *BaseSpider) Close() {
	atomic.StoreUint32(&s.exit, 1)
	s.Logger.Info("==========================spider close====================================")
}

func init() {
	// 配置go使用的cpu数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
