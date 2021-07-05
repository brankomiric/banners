package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/minus5/svckit/log"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (n *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(rw, "Unable to read request body", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, string(b))
}