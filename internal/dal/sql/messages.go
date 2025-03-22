package sql

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type MessagesDAL struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      // Set to current time if it is zero on creating
	Updated   int64          `gorm:"autoUpdateTime:nano"` // Use unix nano seconds as updating time
	Created   int64          `gorm:"autoCreateTime"`      // Use unix seconds as creating time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserId    int64
	MessageId []byte
}

func init() {
	tables = append(tables, &MessagesDAL{})
}

func (dal *Dal) SaveMessage(ctx context.Context, userid int64, messageid []byte) error {
	msgdal := MessagesDAL{
		UserId:    userid,
		MessageId: messageid,
	}
	return dal.db.WithContext(ctx).Create(&msgdal).Error
}

func (dal *Dal) GetMessagesByUserId(ctx context.Context, userid int64) ([]MessagesDAL, error) {
	var msgdal []MessagesDAL
	err := dal.db.
		WithContext(ctx).
		Where("user_id = ?", userid).
		Order("created_at ASC").
		// Limit(limit).
		Find(&msgdal).
		Error

	if err != nil {
		return nil, err
	}
	return msgdal, nil
}
