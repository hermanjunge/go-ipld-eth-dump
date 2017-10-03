package main

import (
	"encoding/binary"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type gethDB struct {
	db *leveldb.DB
}

func gethDBInitStart(path string) *gethDB {
	if path == "" {
		panic("Path to the Geth's DB must be specified (--geth-db-filepath option)")
	}
	db, err := leveldb.OpenFile(path, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		fmt.Println("Corrupta")
		db, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		panic(err)
	}

	return &gethDB{db: db}
}

func (g *gethDB) Stop() {
	g.db.Close()
}

func (g *gethDB) Get(key []byte) ([]byte, error) {
	return g.db.Get(key, nil)
}

func (g *gethDB) GetCanonicalHash(number uint64) []byte {
	headerPrefix := []byte("h")
	numSuffix := []byte("n")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), numSuffix...)
	val, _ := g.db.Get(key, nil)

	return val
}

func (g *gethDB) GetHeaderRLP(hash []byte, number uint64) []byte {
	headerPrefix := []byte("h")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), hash...)

	val, _ := g.db.Get(key, nil)
	return val
}

func (g *gethDB) GetBodyRLP(hash []byte, number uint64) []byte {
	bodyPrefix := []byte("b")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(bodyPrefix, encodedNumber...), hash...)

	val, _ := g.db.Get(key, nil)
	return val
}
