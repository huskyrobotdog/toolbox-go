package database

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/huskyrobotdog/toolbox-go/id"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrAffectedIncorrect = errors.New("rows affected incorrect")

func IsDuplicateEntry(err error) bool {
	if errMySQL, ok := err.(*mysql.MySQLError); ok {
		switch errMySQL.Number {
		case 1062:
			return true
		default:
			return false
		}
	}
	return false
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func CheckResult(result *gorm.DB, matchedRows ...int64) error {
	if result.Error != nil {
		return result.Error
	}
	if len(matchedRows) > 0 {
		if result.RowsAffected != matchedRows[0] {
			return ErrAffectedIncorrect
		}
	}
	return nil
}

func ForUpdate(db *gorm.DB) *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}

func Find[M any](db *gorm.DB, condition ...func(tx *gorm.DB)) ([]M, error) {
	var list []M
	sql := db.Model(new(M))
	if len(condition) > 0 {
		condition[0](sql)
	}
	if err := sql.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func Take[M any](db *gorm.DB, condition ...func(tx *gorm.DB)) (*M, error) {
	var m M
	sql := db.Model(m)
	if len(condition) > 0 {
		condition[0](sql)
	}
	if err := sql.Take(&m).Error; err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func TakeByID[M any](db *gorm.DB, id id.ID) (*M, error) {
	var m M
	if err := db.Where("id=?", id).Take(&m).Error; err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func Create[M any](db *gorm.DB, m *M) error {
	result := db.Create(m)
	if err := CheckResult(result, 1); err != nil {
		return err
	}
	return nil
}

func CreateInBatches[M any](db *gorm.DB, ms []M) error {
	if len(ms) == 0 {
		return nil
	}
	result := db.CreateInBatches(ms, 100)
	if err := CheckResult(result, int64(len(ms))); err != nil {
		return err
	}
	return nil
}

func Delete[M any](db *gorm.DB, condition ...func(tx *gorm.DB)) (int64, error) {
	sql := db
	if len(condition) > 0 {
		condition[0](sql)
	}
	result := sql.Delete(new(M))
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func DeleteByIds[M any](db *gorm.DB, ids ...id.ID) error {
	result := db.Where("id in ?", ids).Delete(new(M))
	if err := CheckResult(result, int64(len(ids))); err != nil {
		return err
	}
	return nil
}
