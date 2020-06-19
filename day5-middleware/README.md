## 中间件

中间件middleware其实和路由映射的HandlerFunc一样， 都是处理方法

##### 增加中间件logger.go
```
package gee

import (
	"log"
	"time"
)

//日志中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()//表示等待执行其他的中间件或用户的Handler

		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
```

##### context.go的变化
```
type Context struct {
	Writer 	http.ResponseWriter
	Req 	*http.Request
	Path 	string
	Method 	string
	Params  map[string]string  
	StatusCode	int

	//中间件
	handlers []HandlerFunc //存放中间件， 方便最后通过上下文来调用
	index int //在Next()中用来计数
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

//返回错误的json响应
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.Json(code, H{"message":err})
}
```

##### gee.go变化
```
//在路由组中添加中间件
func (group *RouterGroup) Use(middleware ...HandlerFunc)  {
	group.middlewares = append(group.middlewares, middleware...)
}


func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	var middlewares []HandlerFunc
	for _, group := range engine.groups{
		if strings.HasPrefix(req.URL.Path, group.prefix) { //根据前缀匹配该路由要用到的所有中间件
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares  //把该路由需要的中间件传给context
	engine.router.handle(c)
}
```

##### router.go的变化
``` 
//处理响应
func (r *router) handle(c *Context)  {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params //为Context的Params赋值

		key := c.Method + "-" + n.pattern
		//r.handlers[key](c)  //根据路由调用对应handler
		c.handlers = append(c.handlers, r.handlers[key])
	}else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s \n", c.Path)
		})
	}
	//通过context调用handlerFunc
	c.Next()
}
```

### c.Next()的调用顺序
```
func A(c *Context) {
    part1
    //中间件上半部分
    c.Next()
    //中间件下半部分
    part2
}
func B(c *Context) {
    part3
    //中间件上半部分
    c.Next()
    //中间件下半部分
    part4
}
```
如果现在有中间件A和B，还有一个路由映射的HandlerFunc
c.handlers是这样的[A, B, HandlerFunc]，c.index初始化为-1。调用c.Next()，调用流程是
```
c.index++，c.index 变为 0
0 < 3，调用 c.handlers[0]，即 A
执行 part1，调用 c.Next()
c.index++，c.index 变为 1
1 < 3，调用 c.handlers[1]，即 B
执行 part3，调用 c.Next()
c.index++，c.index 变为 2
2 < 3，调用 c.handlers[2]，即Handler
Handler 调用完毕，返回到 B 中的 part4，执行 part4
part4 执行完毕，返回到 A 中的 part2，执行 part2
part2 执行完毕，结束
```
最终顺序 part1 -> part3 -> Handler -> part 4 -> part2
