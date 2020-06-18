package main

import (
	"gee"
	"net/http"
)

func main()  {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "hello Gee\n")
	})

	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geehszz"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9090")
}
