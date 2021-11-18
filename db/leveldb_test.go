package db

import (
	"testing"
	"xfsmiddle/common"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func newDb() (*Database, error) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}
func Test_Put(t *testing.T) {
	db, err := newDb()
	if err != nil {
		t.Fatal(err)
		return
	}
	key := "test:"
	value := "testdata"
	if err = db.PutStr(key, value); err != nil {
		t.Fatal(err)
		return
	}
}

func Test_Get(t *testing.T) {
	db, err := newDb()
	if err != nil {
		t.Fatal(err)
		return
	}
	key := "test:"
	value := "testdata"
	if err = db.PutStr(key, value); err != nil {
		t.Fatal(err)
		return
	}

	data, err := db.GetStr(key)
	if err != nil {
		t.Fatal(err)
		return
	}
	if !common.Equal(data, []byte(value)) {
		t.Fatal("Inconsistent acquisition data")
		return
	}
}

func Test_batch(t *testing.T) {
	db, err := newDb()
	if err != nil {
		t.Fatal(err)
		return
	}

	// create db batch object
	db.NewBatch()

	data := make(map[string]string)
	data["test"] = "testdata"
	data["testdata"] = "testdata1"
	data["token"] = "testdata2"

	if err := db.BatchAdd(data); err != nil {
		t.Fatal(err)
		return
	}

	db.BatchDel("token")

	if err := db.BatchCommit(); err != nil {
		t.Fatal(err)
		return
	}
	i := 0
	db.Foreach(func(k string, v []byte) error {
		i++
		return nil
	})
	if i != 2 {
		t.Fatal("The expected data is inconsistent with the actual data")
		return
	}
}
