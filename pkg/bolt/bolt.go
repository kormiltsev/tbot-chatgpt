package bolt

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB(db *bolt.DB) *BoltDB {
	BoltDB := &BoltDB{
		db: db,
	}

	return BoltDB
}

func (db *BoltDB) IsBucketExists(bucketName string) bool {
	tx, err := db.db.Begin(true)
	if err != nil {
		return false
	}
	defer tx.Rollback()

	if bucket := tx.Bucket([]byte(bucketName)); bucket != nil {
		return true
	}

	return false
}

func (db *BoltDB) CreateBucketsIfNotExists(bucketNames []string) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, tbl := range bucketNames {
		bucket := tx.Bucket([]byte(tbl))
		if bucket == nil {
			if _, err := tx.CreateBucketIfNotExists([]byte(tbl)); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (db *BoltDB) Get(bucketName string, key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("empty key")
	}

	var value []byte
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get(key)
		if v != nil {
			value = append(value, v...)
			return nil
		}
		return fmt.Errorf("not found")
	})

	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, fmt.Errorf("not found")
	}

	return value, nil
}

func (db *BoltDB) Put(bucketName string, key, value []byte) error {
	if len(key) == 0 {
		return fmt.Errorf("empty key")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		return b.Put(key, value)
	})
}

func (db *BoltDB) Delete(bucketName string, key []byte) error {
	if _, err := db.Get(bucketName, key); err != nil {
		return err
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete(key)
	})
}

func (db *BoltDB) DeleteWithoutChecking(bucketName string, key []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete(key)
	})
}

func (db *BoltDB) GetAll(bucketName string) ([][]byte, error) {
	var result [][]byte

	return result, db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.ForEach(func(k, v []byte) error {
			// slice can be reused, so copy it:
			dstv := make([]byte, len(v))
			copy(dstv, v)

			result = append(result, dstv)
			return nil
		})
	})
}

func (db *BoltDB) DeleteBucket(bucketName string) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := tx.DeleteBucket([]byte(bucketName)); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *BoltDB) RecreateBucket(bucketName string) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := tx.DeleteBucket([]byte(bucketName)); err != nil {
		return err
	}
	if _, err := tx.CreateBucket([]byte(bucketName)); err != nil {
		return err
	}
	return tx.Commit()
}
