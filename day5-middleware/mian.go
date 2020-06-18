package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

//自定义组件
func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()

		c.Fail(500, "Internal Server Error")
		//fmt.Println("hsz22")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main()  {
	r := gee.New()
	r.Use(gee.Logger())
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee<h1>")
	})

	v2 := r.Group("/v2") //参数为v2时没有执行中间件onlyForV2,应为onlyForV2的v2和/v2不匹配没有加入到c.handlers
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s \n", c.Param("name"), c.Path)
		})
	}
	r.Run(":9090")
}
