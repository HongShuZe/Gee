package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node //动态路由  前缀树
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//把路由分解为字符串切片
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") //以 / 分割pattern

	parts := make([]string, 0)
	for _, val := range vs {
		if val != "" {
			parts = append(parts, val)
			if val[0] == '*' {
				break
			}
		}
	}
	return parts //[hello, :name]
}

//添加路由映射和路由前端树
func (r *router) addRoute(method string, pattern string, handler HandlerFunc)  {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}


//返回节点和 动态路由对应的值的映射
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts{
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

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

