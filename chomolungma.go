package chomolungma

import (
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by orange
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	r      *router
	groups []*RouterGroup // store all groups
}

// New is the constructor of orange.Engine
func New() *Engine {
	engine := &Engine{
		r:      newRouter(),
		groups: make([]*RouterGroup, 0, 8),
	}

	// root group
	rootGroup := &RouterGroup{
		prefix:      "/",
		middlewares: make([]HandlerFunc, 0, 4),
		parent:      nil,
		engine:      engine,
	}
	rootGroup.middlewares = append(rootGroup.middlewares, Recovery())

	engine.RouterGroup = rootGroup
	engine.groups = append(engine.groups, rootGroup)

	return engine
}

func (engine *Engine) Group(prefix string) *RouterGroup {
	group := &RouterGroup{
		prefix:      prefix,
		middlewares: nil,
		parent:      engine.RouterGroup,
		engine:      engine,
	}
	engine.groups = append(engine.groups, group)
	return group
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// implements http.Handler interface
func (engine *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(res, req)
	c.handlers = middlewares
	engine.r.handle(c)
}