package bolt

import (
	"bytes"
	"encoding/gob"
	"os"

	pkgbolt "github.com/kormiltsev/tbot-chatgpt/pkg/bolt"

	boltdb "github.com/boltdb/bolt"
)

type Msg struct {
	Role    string
	Message string
}

type BoltdbCtxInterface interface {
	IsBucketExists(bucketName string) bool
	CreateBucketsIfNotExists(bucketNames []string) error
	Get(bucketName string, key []byte) ([]byte, error)
	Put(bucketName string, key []byte, value []byte) error
	Delete(bucketName string, key []byte) error
	GetAll(bucketName string) ([][]byte, error)
	DeleteBucket(bucketName string) error
	DeleteWithoutChecking(bucketName string, key []byte) error
	RecreateBucket(bucketName string) error
}

func NewBoltDB(fileaddress string) (*boltdb.DB, error) {
	db, err := boltdb.Open(fileaddress, os.FileMode(0o660), nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type ItemBoltDB struct {
	table    string
	itemBolt BoltdbCtxInterface
}

func New(db *boltdb.DB, backet string) (ItemBoltDB, error) {
	newdb := pkgbolt.NewBoltDB(db)
	err := newdb.CreateBucketsIfNotExists([]string{backet})
	if err != nil {
		return ItemBoltDB{}, err
	}
	return ItemBoltDB{backet, newdb}, nil
}

func (db ItemBoltDB) Put(id []byte, role, message string) error {
	msg := Msg{Role: role, Message: message}

	var data bytes.Buffer
	err := gob.NewEncoder(&data).Encode(msg)
	if err != nil {
		return err
	}
	return db.itemBolt.Put(db.table, id[:], data.Bytes())
}

func (db ItemBoltDB) Get(id []byte) (*Msg, error) {
	bytesdata, err := db.itemBolt.Get(db.table, id[:])
	if err != nil {
		return nil, err
	}

	var message Msg
	return &message, gob.NewDecoder(bytes.NewBuffer(bytesdata)).Decode(&message)
}
