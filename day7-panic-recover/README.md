## 错误恢复

defer: panic会导致程序被中止，但是在退出前，会先处理完当前协程上已经defer 的任务，执行完成后再退出。
recover: 可以避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效。

##### recovery.go
```
package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//打印调试的堆栈跟踪
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // 返回调用栈的程序计数器

	var str strings.Builder
	str.WriteString(message + "\n TranceBack:")
	for _, pc := range pcs[:n]{
		fn := runtime.FuncForPC(pc) //获取对应函数
		file, line := fn.FileLine(pc) //获取调用该函数的文件名和行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

//错误处理
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {  //捕获panic
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message)) // 打印错误信息
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()//调用其他的HandlerFunc
	}
}
```

##### gee.go新增了一个方法， 默认使用日志和错误处理中间件
```
//默认使用中间件Logger和Recovery
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
```
