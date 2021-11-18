package xfsmiddle

import (
	"xfsmiddle/db"
)

type TokenManage struct {
	tokenDb db.IDatabase
}

func New(n db.IDatabase) *TokenManage {
	return &TokenManage{
		tokenDb: n,
	}
}

func (n *TokenManage) NewToken() {
	// key := secret.NewKey()

}

func (n *TokenManage) DelToken() {

}

func (n *TokenManage) getToken() {

}
