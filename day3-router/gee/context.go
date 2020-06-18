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
	//响应信息
	StatusCode	int
	//增加参数
	Params  map[string]string  //如 :name 对应的实际参数，就可以使用Param（）来获取， 比较方便
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
	}
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

