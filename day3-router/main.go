package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	//http://localhost:9999/
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello<h1>")
	})

	//http://localhost:9999/hello?name=geektutu
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s \n", c.Query("name"), c.Path)
	})

	//http://localhost:9999/hello/geektutu
	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s \n", c.Param("name"), c.Path)
	})

	//http://localhost:9999/assets/css/geektutu.css
	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.Json(http.StatusOK, gee.H{
			"filepath": c.Param("filepath"),
		})
	})

	r.Run(":9999")
}


