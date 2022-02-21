package chomolungma

import (
	"log"
	"net/http"
	"strings"
)

const capacity = 2 << 4

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node, capacity),
		handlers: make(map[string]HandlerFunc, capacity),
	}
}

func parsePattern(pattern string) []string {
	ss := strings.Split(pattern, "/")
	// make slice (type, len, cap)
	parts := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != "" {
			parts = append(parts, s)
			if s[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
}

// 解析路由, 返回匹配路由的节点和路由参数
func (r *router) resolveRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	matched := root.search(searchParts, 0)
	if matched != nil {
		// 解析参数
		parts := parsePattern(matched.pattern) // 匹配到的可能是子路径
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return matched, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	matched, params := r.resolveRoute(c.Method, c.Path)
	if matched != nil {
		key := c.Method + "-" + matched.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		// need a default HandlerFunc to handle link
		c.handlers = append(c.handlers, func(c *Context) {
			c.JSON(http.StatusNotFound, H{"code": http.StatusNotFound, "message": "404 NOT FOUND"})
		})
	}

	c.Next()
}
