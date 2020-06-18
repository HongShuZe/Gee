package gee

import (
	"net/http"
	"strings"
)
//增加动态路由功能， 就是把路由通过前缀树存储起来， 根据前缀树叶子节点的pattern和method组成的key来找到对应的handlerFunc
type router struct {
	handlers map[string]HandlerFunc
	roots map[string]*node //前端树节点映射，get/post对应不同的前缀树
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//例如pattern为 /hello/:name
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

//添加路由映射和前缀树映射
func (r *router) addRoute(method string, pattern string, handler HandlerFunc)  {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)//添加前缀树映射

	r.handlers[key] = handler//添加路由映射
}


//返回节点和 动态路由对应的实际值的映射params
//如动态路由为/hello/:name，实际为/hello/hsz， params[name]=hsz
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

//获取前缀树所有叶子节点，在test中可以验证路由数量
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}


func (r *router) handle(c *Context)  {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params //为Context的Params赋值

		key := c.Method + "-" + n.pattern
		r.handlers[key](c) //根据路由调用对应handler
	}else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s \n", c.Path)
	}
}

