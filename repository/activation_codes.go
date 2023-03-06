package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"pvg/domain"
)

type acRepository struct {
	db *gorm.DB
}

func NewACRepository(d *gorm.DB) domain.ActivationCodeRepository {
	return &acRepository{
		db: d,
	}
}

func (a *acRepository) GetByUserId(ctx context.Context, id int) (domain.ActivationCodes, error) {
	var (
		res domain.ActivationCodes
		err error
	)

	err = a.db.WithContext(ctx).First(&res, "user_id = ?", id).Error
	if err != nil {
		logrus.Errorf("Activation Code - Repository|err when get AC by user_id, err:%v", err)
		return domain.ActivationCodes{}, err
	}

	return res, nil
}
