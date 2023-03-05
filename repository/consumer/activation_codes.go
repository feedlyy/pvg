package consumer

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pvg/domain/consumer"
)

type acRepository struct {
	db *gorm.DB
}

func NewACRepository(db *gorm.DB) consumer.ACRepository {
	return &acRepository{
		db: db,
	}
}

func (a *acRepository) Insert(ctx context.Context, ac consumer.ActivationCodes) error {
	var err error

	if err = a.db.WithContext(ctx).Create(&ac).Error; err != nil {
		logrus.Errorf("Activation Codes - Repository|err when store AC, err:%v", err)
		return err
	}

	return nil
}
