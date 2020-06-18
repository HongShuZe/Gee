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
			if err := recover(); err != nil { //捕获panic
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message)) // 打印错误信息
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
