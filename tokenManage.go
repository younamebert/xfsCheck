package xfsmiddle

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strconv"
	"time"
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

func (n *TokenManage) NewToken(Group string) (string, error) {
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	result := hex.EncodeToString(h.Sum(nil))
	if err := n.tokenDb.Put([]byte(result), []byte(Group)); err != nil {
		return "", err
	}
	return result, nil
}

func (n *TokenManage) DelToken(token string) error {
	return n.tokenDb.Delete([]byte(token))
}

func (n *TokenManage) GetToken(token string) ([]byte, error) {
	return n.tokenDb.GetStr(token)
}

func (n *TokenManage) SetTokenGroup(token, Group string) error {
	return n.tokenDb.PutStr(token, Group)
}
