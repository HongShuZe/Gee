package gee

import (
	"log"
	"net/http"
)

//定义请求处理方法
//参数Request ，该对象包含了该HTTP请求的所有的信息，比如请求地址、Header和Body等信息；
//参数是ResponseWriter ，利用 ResponseWriter 可以构造针对该请求的响应。
type HandlerFunc func(*Context)


//路由部分
//在golang中有个Handler的概念，一个URL对应一个Handler，在Handler中处理request的具体逻辑，对应关系保存在一个map结构中
type (
	RouterGroup struct {
		prefix		string        // 部分路由
		middlewares []HandlerFunc // 支持中间件，中间应用在分组上
		parent 		*RouterGroup  // 支持嵌套,没怎么用到，可删
		engine 		*Engine       // 所有group共享一个Engine实例，Group通过该指针可以访问router
	}

	Engine struct {
		router *router
		*RouterGroup //继承RouterGroup，将Engine作为最顶层的分组，使Engine拥有RouterGroup所有的能力。
		groups []*RouterGroup //存储所有组，在中间件时会使用到
	}
)

//构造函数
func New() *Engine {
	 engine := &Engine{router: newRouter()}
	 engine.RouterGroup = &RouterGroup{engine: engine}
	 engine.groups = []*RouterGroup{engine.RouterGroup}
	 return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//路由映射表
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc)  {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler) //Group通过*Engine指针可以访问router
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
	c := newContext(w, req)
	engine.router.handle(c)
}





