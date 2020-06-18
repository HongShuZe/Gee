package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//方便构造json数据
type H map[string]interface{}

//封装http.ResponseWriter和*http.Request， 并提供html/json/string等返回类型
//context是很重要的， 可以把很多重要内容封装进去（如 中间件等），把扩展性和复杂性留在内部， 对外简化接口
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


func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
	}
}

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

