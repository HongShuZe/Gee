package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//origin objects 源对象
	Writer http.ResponseWriter
	Req    *http.Request
	//请求信息
	Path   string
	Method string
	Params map[string]string //存放动态路由键值对，方便调用（如 :name 对应的实际参数）
	//响应信息
	StatusCode int

	//中间件
	handlers []HandlerFunc //存放中间件， 方便最后通过上下文来调用
	index    int

	engine   *Engine //使context能通过engine访问html模板
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

//调用中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
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

//html template render
func (c *Context) HTML(code int, name string, data interface{})  {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)

	//ExecuteTemplate: 将指定name的模板解析并应用于data，并将输出写到c.Writer
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}




