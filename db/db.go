package db

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/matobet/shabang/model"
)

type HashDB interface {
	// Write persists the key-value pair
	Write(hash, value model.HashBytes) error

	// Check checks whether a pre-image with given hash is persisted,
	// and if yes returns the pre-image. Otherwise returns nil.
	Check(hash model.HashBytes) (model.HashBytes, error)

	// Close frees all resources held by the DB.
	Close() error
}

type hashDb struct {
	*bolt.DB
}

func Open(path string) (HashDB, error) {
	b, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	b.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(BucketName)
		_, err := tx.CreateBucket(BucketName)
		return err
	})
	return &hashDb{b}, nil
}

var BucketName = []byte("hashes")

func (db *hashDb) Write(hash, value model.HashBytes) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketName)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found", BucketName)
		}

		return bucket.Put(hash, value)
	})
}

func (db *hashDb) Check(hash model.HashBytes) (res model.HashBytes, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketName)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", BucketName)
		}

		res = bucket.Get(hash)
		return nil
	})
	return
}
