package gee

import (
	"log"
	"time"
)

//日志中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()//表示等待执行其他的中间件或用户的Handler

		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
