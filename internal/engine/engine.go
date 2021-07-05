package engine

import (
	"fmt"
	"net/http"
)

type Engine struct {
}

func New() *Engine {
	return &Engine{}
}

func (n *Engine) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "dummy response")
}