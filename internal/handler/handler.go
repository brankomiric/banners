package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/autocorrectoff/banners/internal/dto"
	"github.com/minus5/svckit/log"
)

type Handler struct{
	rpo Repo
}

type Repo interface {
	InsertOrUpdate(match *dto.Match) error
	FindByIdIn(ids []int) ([]dto.Match, error)
}

func New(repo Repo) *Handler {
	return &Handler{
		rpo: repo,
	}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	handleError(rw, err, "Unable to read request body", http.StatusBadRequest)
	
	parsedBody := &struct { 
		MatchIds []int
	}{
		MatchIds: make([]int, 0),
	}
	err = json.Unmarshal(b, parsedBody)
	handleError(rw, err, "Unable to unmarshal request body", http.StatusBadRequest)
	
	matches, dbErr := h.rpo.FindByIdIn(parsedBody.MatchIds)
	handleError(rw, dbErr, "Unable to marshal response", http.StatusInternalServerError)

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(rw)
	encodingErr := e.Encode(matches)
	handleError(rw, encodingErr, "Unable to marshal response", http.StatusInternalServerError)
}

func handleError(rw http.ResponseWriter, err error, errorMessage string, errorStatus int) {
	if err != nil {
		log.Error(err)
		http.Error(rw, errorMessage, errorStatus)
		return
	}
}