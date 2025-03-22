package sql

import (
	"gorm.io/gorm"
)

// New creates DAL abstraction over sql.DB type
func New(db *gorm.DB) *Dal {
	return &Dal{db: db}
}

// Dal is a database abstraction layer
type Dal struct {
	db *gorm.DB
}

var tables []interface{} = make([]interface{}, 0)

func (dal *Dal) Migrate() error {
	for _, model := range tables {
		err := dal.db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}
