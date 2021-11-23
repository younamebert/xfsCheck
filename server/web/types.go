package web

type Empty struct{}

type NewTokenArgs struct {
	Group string `json:"group"`
}

type GetGroupByTokenArgs struct {
	Token string `json:"token"`
}

type PutTokenGroupArgs struct {
	Token string `json:"token"`
	Group string `json:"group"`
}
