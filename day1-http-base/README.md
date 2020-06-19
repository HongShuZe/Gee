
一、我们先看一下标准库net/http处理请求的方式
```
func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

```
mian方法里面的 http.HandleFunc 就是实现路由和相应的处理方法的映射。
http.ListenAndServe是用来启动web服务的的第一个参数表示localhost:8000的端口在监听，第二个参数为nil表示使用标准库的实例处理， 是 我们基于net/http的web框架的入口。

二、Handler接口
```
package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
```
ListenAndServe的第二个参数代表处理所有http请求的实例，是我们基于net/http标准库实现web框架的接口，
Handler是一个接口， 需要实现方法ServeHTTP。 只要传入实现该接口的实例， 所有http请求都交给该实例处理。


三、代码实现
我们可以建一个结构体Engine，里面保存路由映射表（即每个路由对有应不同的处理方法HandlerFunc），并实现ServeHTTP作为框架入口， 还有增加路由的方法

##### 架构雏形
```
gee/
  |--gee.go
  |--go.mod
main.go
go.mod
```

##### go.mod
replace gee => ./gee是将gee指向./gee, 即可引用相对路径的package
```
module Gee/day1-http-base

go 1.13

require gee v0.0.0

replace gee => ./gee

```

##### gee.go
```
package gee

import (
	"fmt"
	"net/http"
)

//定义请求处理方法
//参数Request ，该对象包含了该HTTP请求的所有的信息，比如请求地址、Header和Body等信息；
//参数ResponseWriter ，利用 ResponseWriter 可以构造针对该请求的响应。
type HandlerFunc func(http.ResponseWriter, *http.Request)



//在golang中有个Handler的概念，一个URL对应一个Handler，在Handler中处理request的具体逻辑，对应关系保存在一个map结构中
type Engine struct {
	router map[string]HandlerFunc
}


//构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

//路由映射表
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc)  {
	key := method + "-" + pattern
	engine.router[key] = handler
}

//GET请求方法
func (engine *Engine) GET(pattern string, handler HandlerFunc)  {
	engine.addRouter("GET", pattern, handler)
}

//POST请求方法
func (engine *Engine) POST(pattern string, handler HandlerFunc)  {
	engine.addRouter("POST", pattern, handler)
}


//定义http服务器启动方法
//Engine实现了ServeHTTP， 在ListenAndServe 中传入Engine的实例， 即将将http请求交由Engine的实例处理
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

/*type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}*/
//实现ServeHTTP方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	key := req.Method + "-" + req.URL.Path

	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	}else {
		fmt.Fprintf(w, "404 NOT FOUND %s\n", req.URL)
	}
}
```
##### mian.go
```
package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		//格式化并输出到 io.Writer接口类型的变量（只要实现了Write方法）
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header{
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")
}

```

三、 部分解析
HandlerFunc是定义请求处理方法，参数Request ，该对象包含了该HTTP请求的所有的信息，比如请求地址、Header和Body等信息；参数ResponseWriter ，利用 ResponseWriter 可以构造针对该请求的响应。

Engine 里面有一个router的路由映射表， 用来保存url对应的HandlerFunc。
而且Engine实现了方法ServeHTTP，作为框架的入口

GET和POST方法就是将路由和处理方法注册路由表中。ServeHTTP就可以根据路由表调用对应的HandlerFunc。


[此gee框架是参考geektutu](https://geektutu.com/post/gee.html)

