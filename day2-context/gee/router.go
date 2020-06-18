package gee

import "net/http"


//把路router由映射表提取出来, 方便增强
type router struct {
	Handlers map[string]HandlerFunc
}

//构造函数
func newRouter() *router {
	return &router{Handlers: make(map[string]HandlerFunc)}
}


func (r *router) addRouter(method string, pattern string, handler HandlerFunc)  {
	key := method + "-" + pattern
	r.Handlers[key] = handler
}

//由router处理根据key查找对应的handlerFuc方法
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.Handlers[key]; ok {
		handler(c)
	}else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}


