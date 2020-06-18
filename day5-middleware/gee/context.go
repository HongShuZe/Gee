package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//origin objects 源对象
	Writer 	http.ResponseWriter
	Req 	*http.Request
	//请求信息
	Path 	string
	Method 	string
	Params  map[string]string  //存放动态路由键值对，方便调用（如 :name 对应的实际参数）
	//响应信息
	StatusCode	int

	//中间件
	handlers []HandlerFunc //存放中间件， 方便最后通过上下文来调用
	index int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		index:      -1,
	}
}

//index是记录当前执行到第几个中间件，当在中间件中调用Next方法时，控制权交给了下一个中间件，直到调用到最后一个中间件，然后再从后往前，调用每个中间件在Next方法之后定义的部分。
//调用中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	//使用for的原因：不是所有的handlerFunc都在内部调用c.Next(),即手工调用 Next()，防止有handlerFunc没有被调用
	//手工调用 Next()一般用在请求前后各实现一些行为。如果中间件只作用于请求前，可以省略调用Next()。
	for ;c.index < s ; c.index++ {
		c.handlers[c.index](c)
	}
}


func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.Json(code, H{"message":err})
}


//获取动态路由对应参数， 如:name 对应的参数
func (c *Context) Param(key string) string {
	value , _ := c.Params[key]
	return value
}

func (c *Context) PostFrom(key string) string {
	return c.Req.FormValue(key) //解析url参数
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)  //解析url参数
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

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




