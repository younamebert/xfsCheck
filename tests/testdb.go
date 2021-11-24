package test

import (
	db "xfsmiddle/db"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func NewDb() (db.IDatabase, error) {
	storage, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}

	return &db.Database{
		Db: storage,
	}, nil
}
