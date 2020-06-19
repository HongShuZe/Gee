## 路由分组

路由分组是web框架的基本功能之一，
如果没有分组， 就需要对每一个路由进行控制 
在真实的业务场景中， 一般需要某一组路由进行相似的处理， 
而且中间件就是应用在分组上的，作用在多个路由上

##### gee.go变化
```
package gee

type (
    //新增路由分组结构体
	RouterGroup struct {
		prefix		string        // 部分路由
		middlewares []HandlerFunc // 支持中间件，中间应用在分组上
		parent 		*RouterGroup  // 支持嵌套, 没怎么用到，可删
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

//路由分组方法
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
	//pattern变为分组的部分路由+路由
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
```
