package db

import (
	"go.etcd.io/bbolt"
	"os"
)

func InitKV(path string, buckets ...string) (*bbolt.DB, error) {
	db, err := bbolt.Open(path, os.ModePerm, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range buckets {
			if _, err = tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
