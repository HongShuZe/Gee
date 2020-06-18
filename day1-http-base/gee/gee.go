package gee

import (
	"fmt"
	"net/http"
)

//定义请求处理方法
//参数Request ，该对象包含了该HTTP请求的所有的信息，比如请求地址、Header和Body等信息；
//参数ResponseWriter ，利用 ResponseWriter 可以构造针对该请求的响应。
type HandlerFunc func(http.ResponseWriter, *http.Request)


//路由部分
//在golang中有个Handler的概念，一个URL对应一个Handler，在Handler中处理request的具体逻辑，对应关系保存在一个map结构中
type Engine struct {
	//路由映射表
	router map[string]HandlerFunc
}


//构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

//添加到路由映射表
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





