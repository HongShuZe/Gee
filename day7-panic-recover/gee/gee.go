package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
		router 		*router
		*RouterGroup
		groups 		[]*RouterGroup //存储所有组

		htmlTemplates *template.Template //对html渲染 (生成安全的html片段)
		funcMap      template.FuncMap //对html渲染 (定义从名称到函数的映射)
		//htmlTemplates将所有的模板加载进内存，funcMap是所有的自定义模板渲染函数。
	}
)

//构造函数
func New() *Engine {
	 engine := &Engine{router: newRouter()}
	 engine.RouterGroup = &RouterGroup{engine: engine}
	 engine.groups = []*RouterGroup{engine.RouterGroup}
	 return engine
}

//默认使用中间件Logger和Recovery
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
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
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

//创建静态handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs)) //StripPrefix返回服务HTTP请求的处理程序 通过从请求URL的路径中删除给定的前缀并调用处理程序h.

	return func(c *Context) {
		file := c.Param("filepath")

		if _, err := fs.Open(file); err != nil{
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}


func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))

	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)//添加路由映射
}

//设置自定义渲染函数funcMap
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//加载模板的方法
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}


