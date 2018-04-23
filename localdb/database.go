// Copyright 2017. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package localdb

import (
	"github.com/boxproject/voucher/config"
	verrors "github.com/boxproject/voucher/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	logger "github.com/alecthomas/log4go"
)

type Database interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	GetPrifix(keyPrefix []byte) (map[string]string, error)
	Close() error
}

func NewLVDBDatabase(filePath string, cache int, openFiles int) (*leveldb.DB, error) {
	if cache < 16 {
		cache = 16
	}
	if openFiles < 16 {
		openFiles = 16
	}

	db, err := leveldb.OpenFile(filePath, &opt.Options{
		OpenFilesCacheCapacity: openFiles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB,
		Filter:                 filter.NewBloomFilter(10),
	})

	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(filePath, nil)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewDatabase(cfg *config.DatabaseConfig) (Database, error) {
	db, err := NewLVDBDatabase(cfg.Filepath, cfg.Cache, cfg.Openfiles)
	if err != nil {
		return nil, err
	}

	return &table{
		db:     db,
		prefix: cfg.Prefix,
	}, nil
}

type table struct {
	db     *leveldb.DB
	prefix string
}

func (dt *table) Put(key []byte, value []byte) error {
	return dt.db.Put(append([]byte(dt.prefix), key...), value, nil)
}

func (dt *table) Has(key []byte) (ret bool, err error) {
	ret, err = dt.db.Has(append([]byte(dt.prefix), key...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return ret, verrors.NoDataErr
		}
		return
	}

	return ret, nil
}

func (dt *table) Get(key []byte) ([]byte, error) {
	data, err := dt.db.Get(append([]byte(dt.prefix), key...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, verrors.NoDataErr
		}
		return nil, err
	}
	return data, nil
}

//查询前缀 TODO iter 使用字节组赋值有问题
func (dt *table) GetPrifix(keyPrefix []byte) (map[string]string, error) {
	var resMap map[string]string = make(map[string]string)
	iter := dt.db.NewIterator(util.BytesPrefix(append([]byte(dt.prefix), keyPrefix...)), nil)
	if iter.Error() == leveldb.ErrNotFound {
		return nil, verrors.NoDataErr
	}
	if iter.Error() != nil {
		logger.Error("get prifix error")
		return nil, iter.Error()
	}
	for iter.Next() {
		resMap[string(iter.Key())] = string(iter.Value())
	}

	iter.Release()

	return resMap, nil
}

func (dt *table) Delete(key []byte) error {
	err := dt.db.Delete(append([]byte(dt.prefix), key...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return verrors.NoDataErr
		}

		return err
	}

	return nil
}

func (dt *table) Close() error {
	return dt.db.Close()
}
