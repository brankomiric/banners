package engine

import "github.com/autocorrectoff/banners/internal/dto"

type Engine struct {
	rpo Repo
}

type Repo interface {
	InsertOrUpdate(match *dto.Match) error
}

func New(repo Repo) *Engine {
	return &Engine{
		rpo: repo,
	}
}

func (n *Engine) HandleMatch(req *dto.MatchRequest) error {
	var oddsList []dto.DisplayOdds
	for _, odd := range(req.Offers[0].Odds) {
		o := dto.DisplayOdds{
			Odd: odd.Odd,
			Name: odd.Title,
			Position: odd.Order,
		}
		oddsList = append(oddsList, o)
	}
	match := &dto.Match{
		BaseId:      req.BaseId,
		DisplayOdds: oddsList,
	}
	err := n.rpo.InsertOrUpdate(match)
	if err != nil {
		return err
	}
	return nil
}