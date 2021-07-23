/*
 * @Author: xiangcai
 * @Date: 2021-07-23 17:11:54
 * @LastEditors: xiangcai
 * @LastEditTime: 2021-07-23 17:47:35
 * @Description: file content
 */

 package gspider

 import (
	"time"
	 "sync/atomic"
 )


 type CounterExtension struct {
	curReqCounter   int64
	curRespCounter  int64
	reqCounter      int64
	respCounter     int64
 }

 // 给请求计数器追加1
func (c *CounterExtension) recordReq(*Request) {
	atomic.AddInt64(&c.curReqCounter, 1)
	atomic.AddInt64(&c.reqCounter, 1)
}

// 给请求计数器追加1
func (c *CounterExtension) recordResp(*Response) {
	atomic.AddInt64(&c.curRespCounter, 1)
	atomic.AddInt64(&c.respCounter, 1)
}

func (c *CounterExtension) Run(spider *BaseSpider) {
	spider.OnRequest(c.recordReq)
	spider.OnResponse(c.recordResp)
	ticker := time.NewTicker(time.Minute)
	for {
		if atomic.LoadUint32(&spider.exit) != 0 {
			break
		}
		select {
		case <-ticker.C:
			spider.Logger.Infof(
				"----------------------crawled (%d/%d), (%d/%d)/min------------------------",
				atomic.LoadInt64(&c.reqCounter),
				atomic.LoadInt64(&c.respCounter),
				atomic.SwapInt64(&c.curReqCounter, 0),
				atomic.SwapInt64(&c.curRespCounter, 0),
			)
		}
	}
}


func NewCounterExtension() Extension {
	return &CounterExtension{}
}
