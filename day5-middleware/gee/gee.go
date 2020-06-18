package gee

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//定义请求处理方法
//参数Request ，该对象包含了该HTTP请求的所有的信息，比如请求地址、Header和Body等信息；
//参数是ResponseWriter ，利用 ResponseWriter 可以构造针对该请求的响应。
type HandlerFunc func(*Context)

//路由部分
//在golang中有个Handler的概念，一个URL对应一个Handler，在Handler中处理request的具体逻辑，对应关系保存在一个map结构中


type (
	RouterGroup struct {
		prefix		string
		middlewares []HandlerFunc //支持中间件（中间件就是 自定义/默认定义处理程序（HandlerFunc）），
		parent 		*RouterGroup  //支持嵌套
		engine 		*Engine       //所有group共享一个Engine实例
	}

	Engine struct {
		router *router
		*RouterGroup
		groups []*RouterGroup //存储所有组
	}
)

//构造函数
func New() *Engine {
	 engine := &Engine{router: newRouter()}
	 engine.RouterGroup = &RouterGroup{engine: engine}
	 engine.groups = []*RouterGroup{engine.RouterGroup}
	 return engine
}

//在路由组中添加中间件
func (group *RouterGroup) Use(middleware ...HandlerFunc)  {
	group.middlewares = append(group.middlewares, middleware...)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	fmt.Println(newGroup.prefix)
	return newGroup
}

//路由映射表
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc)  {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//GET请求方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc)  {
	group.addRouter("GET", pattern, handler)
}

//POST请求方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc)  {
	group.addRouter("POST", pattern, handler)
}


//定义http服务器启动方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

/*type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}*/
//实现ServeHTTP方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	var middlewares []HandlerFunc
	for _, group := range engine.groups{
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares  //为context的中间件数组赋值
	engine.router.handle(c)
}





