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