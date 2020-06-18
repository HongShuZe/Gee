package main

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

//在templates文件夹的custom_func.tmpl有指定使用该模板
func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf( "hsz-%d-%02d-%02d", year, month, day)
}

func main()  {
	r := gee.New()
	r.Use(gee.Logger())

	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})

	r.LoadHTMLGlob("templates/*") //加载html模板，输入参数： html模板路径

	r.Static("/assets", "./static") //静态文件映射

	stu1 := &student{
		Name: "hsz",
		Age:  22,
	}
	stu2 := &student{
		Name: "hszz",
		Age:  20,
	}

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title": "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/data", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now": time.Date(2020, 6, 15, 0,0,0,0,time.UTC),
		})
	})

	r.Run(":9090")
}
