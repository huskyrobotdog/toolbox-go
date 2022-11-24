package database

import (
	"errors"

	"gorm.io/gorm"
)

var ErrAffectedIncorrect = errors.New("rows affected incorrect")

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func Create[M any](db *gorm.DB, m *M) error {
	result := db.Create(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrAffectedIncorrect
	}
	return nil
}

func CreateInBatches[M any](db *gorm.DB, ms []M) error {
	if len(ms) == 0 {
		return nil
	}
	result := db.CreateInBatches(ms, 100)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != int64(len(ms)) {
		return ErrAffectedIncorrect
	}
	return nil
}
