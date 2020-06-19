##前缀树实现动态路由

前缀树是N叉树的一种特殊形式。一个前缀树可以用来存储字符串，前缀树的每一个节点代表一个字符串（前缀）。每一个节点会有多个子节点，通往不同子节点的路径上有着不同的字符。子节点代表的字符串是由节点本身的原始字符串，以及通往该子节点路径上所有的字符组成的。
![前缀树示例图](https://img2018.cnblogs.com/blog/1519578/201907/1519578-20190724132134884-1903210243.png)

http的路径是根据/来分隔的， 我们可以通过把路由根据/分割然后存进前缀树

####trie.go
```
package gee

import (
	"fmt"
	"strings"
)

//前缀树
type node struct {
	pattern string //待匹配路由， 如 /p/:lang （在最后一个节点保存完整路由，在查询时可判断该节点是否正确）
	part 	string //路由中的一部分， 如 /p, :lang
	children []*node //子节点
	isWild 	bool //路由是否精确匹配， part含有:和*时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

//添加节点
//height为parts数组的坐标，从零开始
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { //parts为空时或为结束递归条件（len（parts）=1， height=0，执行一次后len（parts）= height+1）
		n.pattern = pattern
		return
	}
	
	part := parts[height]
	child := n.matchChild(part)
	if child == nil { //没有该节点就新建
		child = &node{
			part:     part,
			isWild:   part[0] == ':' || part[0] == '*',
		}

		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//height为parts数组的坐标，从零开始
//strings.HasPrefix测试字符串是否以prefix开头
//查询节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //parts为空时 或 为结束递归条件 或 parts为*
		if n.pattern == "" { //验证该节点是否正确
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children{
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

//获取pattern不为空的节点即叶子节点
func (n *node) travel(list *([]*node))  {
	if n.pattern != "" {
		*list = append(*list, n)
	}

	for _, child := range n.children{
		child .travel(list)
	}
}

//第一个匹配成功的节点， 用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild { //child.isWild为true时即该节点有:/*
			return child
		}
	}

	return nil
}

//所有匹配成功的节点， 用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children{
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
```

```
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
		c.Params = params  //为Context的Params赋值

		key := c.Method + "-" + n.pattern
		r.handlers[key](c)  //根据路由调用对应handler
	}else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s \n", c.Path)
	}
}
```

```

```
