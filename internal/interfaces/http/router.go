// Package http defines the HTTP router interface and implementation
package http

import (
	"github.com/gin-gonic/gin"
)

// RouterRegister is the interface a module must implement
type RouterRegister interface {
	Register(router *gin.RouterGroup)
}

type Router struct {
	engine *gin.Engine
}

func NewRouter(engine *gin.Engine) *Router {
	return &Router{engine: engine}
}

// RegisterModules takes a slice of modules and lets them register themselves
func (r *Router) RegisterModules(v1Group *gin.RouterGroup, modules ...RouterRegister) {
	for _, m := range modules {
		m.Register(v1Group)
	}
}
