package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

//const KEY_PREFIX = "dfc:cache:"
const KEY_PREFIX = ""

//资料存储层
type ResourceDB struct {
	*leveldb.DB
}

func OpenResourceDB(filepath string) (*ResourceDB, error) {
	db, err := leveldb.OpenFile(filepath, nil)
	if err != nil {
		return nil, err
	}
	resourceDB := &ResourceDB{DB: db}
	return resourceDB, nil
}

func (self *ResourceDB) Get(key []byte) ([]byte, error) {
	return self.DB.Get(key, nil)
}

func (self *ResourceDB) Set(key []byte, value []byte) error {
	return self.DB.Put(key, value, nil)
}

func (self *ResourceDB) Remove(key []byte) error {
	return self.DB.Delete(key, nil)
}

func (self *ResourceDB) Update(key []byte, value []byte) error {
	batch := new(leveldb.Batch)
	batch.Delete(key)
	batch.Put(key, value)
	return self.DB.Write(batch, nil)
}

func (self *ResourceDB) Keys() []string {
	cumul := []string{}
	iter := self.DB.NewIterator(nil, nil)
	for iter.Next() {
		cumul = append(cumul, string(iter.Key()))
	}
	iter.Release()
	return cumul
}

func (self *ResourceDB) Close() error {
	if self.DB != nil {
		return self.DB.Close()
	}
	return nil
}
