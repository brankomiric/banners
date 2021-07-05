package dto

type MatchRequest struct {
	Source int `json:"izvor"`
	SourceId int `json:"izvorId"`
	BaseId int `json:"baseId"`
	SportId int `json:"sportId"`
	Home string `json:"domacin"`
	Aways string `json:"gost"`
	League string `json:"liga"`
	Region string `json:"region"`
	Time string `json:"vrijeme"`
	Status int `json:"stanje"`
	ChangedAt string `json:"changedAt"` 
	HasStat bool `json:"imaStatistiku"`
	SportRadarId int `json:"sportradarId"`
	Offers []Offer `json:"ponude"`
}

type Offer struct {
	SourceId int `json:"izvorId"`
	DbId int `json:"dbId"`
	Title string `json:"naziv"`
	SportRazradaId int `json:"sportRazradaId"`
	Status int `json:"stanje"`
	Number string `json:"broj"`
	ChangedAt int `json:"changedAt"` 
	Odds []Odd `json:"tecajevi"`
	Combination *Combination `json:"kombiniranje"`
	Handicap string `json:"hendikep"`
}

type Odd struct {
	SourceId int `json:"izvorId"`
	DbId int `json:"dbId"`
	Odd float32 `json:"tecaj"`
	Title string `json:"naziv"`
	Order int `json:"poredak"`
}

type Combination struct {
	DependencyAttr string `json:"dependencyAttr"`
	MinGames int `json:"minUtakmica"`
}

type Match struct {
	BaseId int `bson:"_id"`
	DisplayOdds []DisplayOdds `bson:"odds"`
}

type DisplayOdds struct {
	Odd float32
	Name string
	Position int
}