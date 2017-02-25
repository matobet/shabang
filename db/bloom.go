package db

import (
	"github.com/matobet/shabang/model"
	"github.com/willf/bloom"
)

type bloomAdapter struct {
	HashDB
	filter *bloom.BloomFilter
}

// Bloomify returns an adapter around given HashDB implementation that first checks
// whether the hash could possibly be contained in the DB and only then forwards
// the `Check` request to the underlying implementation.
func Bloomify(db HashDB, m, n uint) HashDB {
	return &bloomAdapter{db, bloom.New(m, n)}
}

func (bloom *bloomAdapter) Write(hash, value model.HashBytes) error {
	bloom.filter.Add(hash)
	return bloom.HashDB.Write(hash, value)
}

func (bloom *bloomAdapter) Check(hash model.HashBytes) (model.HashBytes, error) {
	if !bloom.filter.Test(hash) {
		return nil, nil
	}
	return bloom.HashDB.Check(hash)
}
