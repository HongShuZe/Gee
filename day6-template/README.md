```
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
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
	group.GET(urlPattern, handler)
}
```

在createStaticHandler第二个参数类型为http.FileSystem， 我们传入http.Dir(root)， 是因为Dir实现了FileSystem接口,看源码包
```
type Dir string

func (d Dir) Open(name string) (File, error) {
    ....
}

type FileSystem interface {
    Open(name string) (File, error)
}
```
