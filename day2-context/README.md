##Context上下文

context是在框架起到很大的作用，可以把很多重要内容封装进去（如中间件等， 请求信息等），然后通过访问Context就可以方便的拿到需要的数据。把扩展性和复杂性留在Context内部， 对外简化接口。

我们把http.ResponseWriter和*http.Request封装进Context， 并提供html/json/string等返回类型的方法

####context.go
#####增加结构体Context
```
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//方便构造json数据
type H map[string]interface{}

type Context struct {
	//origin objects 源对象
	Writer 	http.ResponseWriter
	Req 	*http.Request
	//请求信息
	Path 	string
	Method 	string
	//响应信息
	StatusCode	int
}

```

#####新增方法
有了context， 就可以从中方便的取得想要的数据或改变数据
```
//从url解析参数的方法
func (c *Context) PostFrom(key string) string {
	return c.Req.FormValue(key) //解析url参数
}

//从url解析参数的方法
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)  //解析url参数
}

//设置StatusCode
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//设置数据头信息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//各种返回类型

func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "test/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

func (c *Context) Json(code int, obj interface{})  {
	c.SetHeader("Content", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte)  {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string)  {
	c.SetHeader("Content", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
```
要注意http.ResponseWriter的写入顺序
```
c.Writer.Header().Set("Content-type", "application/text")
c.Writer.WriteHeader(200)
c.Writer.Write([]byte(resp))
```

###router.go
把路由映射表提取出来，把和router处理有关的方法提取出来，方便增强，简化gee.go
```
package gee

import "net/http"

type router struct {
	Handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{Handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRouter(method string, pattern string, handler HandlerFunc)  {
	key := method + "-" + pattern
	r.Handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.Handlers[key]; ok {
		handler(c)
	}else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
```

###gee.go
将HandlerFunc的参数简化为Context， 需要使用路由映射表的通过访问router.go
```
package gee

import (
	"log"
	"net/http"
)

//定义请求处理方法
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

//构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

//添加路由映射
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc)  {
	log.Printf("Rount %4s - %s", method, pattern)
	engine.router.addRouter(method, pattern, handler)
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
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//实现ServeHTTP方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	c := newContext(w, req)
	engine.router.handle(c)
}
```


