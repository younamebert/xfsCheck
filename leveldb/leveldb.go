package leveldb

import (
	"xfstoken/logs"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type IDatabase interface {
	Put(key, value []byte) error
	NewBatch()
	BatchAdd(bat map[string]string) error
	BatchDel(key string)
	BatchCommit() error
	PutStr(key, value string) error
	Get(key []byte) ([]byte, error)
	GetStr(key string) ([]byte, error)
	Delete(key []byte) error
	Foreach(fn func(k string, v []byte) error) error
	PrefixForeach(prefix string, fn func(k string, v []byte) error) error
}

type Database struct {
	db    *leveldb.DB
	log   logs.ILogger
	batch *leveldb.Batch
}

func New(pathname string) (*Database, error) {
	db, err := leveldb.OpenFile(pathname, nil)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:  db,
		log: logs.NewLogger("database"),
	}, nil
}

func (db *Database) Put(key, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *Database) NewBatch() {
	db.batch = new(leveldb.Batch)
}

func (db *Database) BatchAdd(bat map[string]string) error {
	for key, val := range bat {
		db.batch.Put([]byte(key), []byte(val))
	}
	return nil
}

func (db *Database) BatchDel(key string) {
	db.batch.Delete([]byte(key))
}

func (db *Database) BatchCommit() error {
	err := db.db.Write(db.batch, nil)
	if err != nil {
		return err
	}
	db.batch.Reset()
	return nil
}

func (db *Database) PutStr(key, value string) error {
	return db.Put([]byte(key), []byte(value))
}

func (db *Database) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *Database) GetStr(key string) ([]byte, error) {
	return db.Get([]byte(key))
}

func (db *Database) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *Database) Foreach(fn func(k string, v []byte) error) error {
	iter := db.db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		if err := fn(string(key), value); err != nil {
			return err
		}
	}
	iter.Release()
	return iter.Error()
}

func (db *Database) PrefixForeach(prefix string, fn func(k string, v []byte) error) error {
	iter := db.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		if err := fn(string(key), value); err != nil {
			return err
		}
	}
	iter.Release()
	return iter.Error()
}
