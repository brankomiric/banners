package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/autocorrectoff/banners/internal/dto"
	"github.com/minus5/svckit/log"
)

const (
	// TODO: replace with actual url
	configUrl = "http://localhost:9000/promo_banners"
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
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	config, err := fetchConfig()
	handleError(rw, err, "Unable to fetch config", http.StatusInternalServerError)
	var matchIds []int
	for _, game := range config.Promo.Games {
		id, convErr := strconv.Atoi(game.Id)
		if convErr != nil {
			log.Errorf("Cannot conver id %s to integer", game.Id)
		}
		matchIds = append(matchIds, id)
	}

	matches, dbErr := h.rpo.FindByIdIn(matchIds)
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

func fetchConfig() (*dto.BannerConfig, error) {
	resp, err := http.Get(configUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	config := &dto.BannerConfig{}
	err = json.Unmarshal(body, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}