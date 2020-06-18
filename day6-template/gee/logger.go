package gee

import (
	"log"
	"time"
)

//日志中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		//fmt.Println("hsz3")
		c.Next()
		//fmt.Println("hsz33")
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
