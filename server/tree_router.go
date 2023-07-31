package server

import (
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// node ...
// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 路径参数匹配：形式 :param_name
// 3. 通配符匹配：*
// 这是不回溯匹配
type node struct {
	path string
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 通配符 * 表达的节点，任意匹配
	starChild *node

	paramChild *node
}

func (tr *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := tr.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	segs := strings.Split(strings.Trim(path, "/"), "/")
	mi := &matchInfo{}
	for _, seg := range segs {
		var matchParam bool
		root, matchParam, ok = root.childOf(seg)
		if !ok {
			return nil, false
		}
		if matchParam {
			mi.addValue(root.path[1:], seg)
		}
	}
	mi.n = root
	return mi, true
}

func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return res, false, ok
}

//goland:noinspection GoMixedReceiverTypes
func (tr router) addRouter(method string, path string, handler HandleFunc) {
	root, ok := tr.trees[method]
	// 这是一个全新的 HTTP 方法，创建根节点
	if !ok {
		// 创建根节点
		root = &node{path: "/"}
		tr.trees[method] = root
	}
	segs := strings.Split(path[1:], "/")
	// 开始一段段处理
	for _, s := range segs {
		if s == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOrCreate(s)
	}
	root.handler = handler
}

func (n *node) findMathChild(path string) (*node, bool) {
	var cardNode *node
	for _, item := range n.children {
		if item.path == path && item.path != "*" {
			return item, true
		}
		if item.path == "*" {
			cardNode = item
		}
	}
	return cardNode, cardNode != nil
}

func (n *node) childOrCreate(path string) *node {
	if path == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.starChild == nil {
			n.starChild = &node{path: path}
		}
		return n.starChild
	}
	if path[0] == ':' {
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.paramChild != nil {
			if n.paramChild.path != path {
				panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
			}
		} else {
			n.paramChild = &node{path: path}
		}
		return n.paramChild
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}
