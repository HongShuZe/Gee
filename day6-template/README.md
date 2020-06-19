##静态文件

##### gee.go
```
//创建静态handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs)) //StripPrefix返回服务HTTP请求的handlerFunc 通过从请求URL的路径中删除给定的前缀并调用处理程序h.
	                                                                  //http.FileServer用来处理访问本地"/tmp"文件夹的HTTP请求
	return func(c *Context) {
		file := c.Param("filepath")//获取*filepath对应的路由， 即为文件最终路径

		if _, err := fs.Open(file); err != nil{
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//设置静态文件映射
//用户可以把磁盘上的文件夹root映射到relativePath
//如路由规则为/relativePath/*filepath，访问localhost:9090/relativePath/hsz 即为本地root/hsz
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root)) //利用本地tmp目录实现一个文件系统

	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)//添加到路由映射表
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
## HTML模板渲染
##### gee.go
```
type (
	
    ...

	Engine struct {
		router 		*router
		*RouterGroup
		groups 		[]*RouterGroup  
		
		htmlTemplates *template.Template  //对html渲染 (生成安全的html片段)
		funcMap      template.FuncMap  //对html渲染 (定义从名称到函数的映射)
		//htmlTemplates将所有的模板加载进内存，funcMap是所有的自定义模板渲染函数。
	}
)

//添加自定义渲染函数funcMap
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//加载模板的方法
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	...

	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine //给context的engine赋值，方便context通过engine访问html模板
	engine.router.handle(c)
}
```
##### context.go
```
type Context struct {
	...

	engine   *Engine //新增*Engine, 使context能通过engine访问html模板
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
```

### 最终目录结构
```
---gee/
---static/
   |---css/
        |---geektutu.css
   |---file1.txt
---templates/
   |---arr.tmpl
   |---css.tmpl
   |---custom_func.tmpl
---main.go

static放着静态文件， templates放着html模板
```


