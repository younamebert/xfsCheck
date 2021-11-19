package web

type NewTokenArgs struct {
	Group string `json:"group"`
}

type DelTokenArgs struct {
	Token string `json:"token"`
}

type GetGroupByTokenArgs struct {
	Token string `json:"token"`
}

type PutTokenGroupArgs struct {
	Token string `json:"token"`
	Group string `json:"group"`
}
