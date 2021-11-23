package sub

type empty struct{}

type newTokenArgs struct {
	group string `json:"group"`
}

type delTokenArgs struct {
	token string `json:"token"`
}

type setTokenArgs struct {
	token string `json:"token"`
	group string `json:"group"`
}

type getTokenGroupArgs struct {
	token string `json:"token"`
}
