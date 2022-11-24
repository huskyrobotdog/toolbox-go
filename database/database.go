package database

import (
	"fmt"
	"time"

	"github.com/huskyrobotdog/toolbox-go/inner"
	"github.com/huskyrobotdog/toolbox-go/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

var Instance *gorm.DB

type Config struct {
	Address     string `json:"address"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	Config      string `json:"config"`
	MaxIdleConn int    `json:"maxIdleConn"`
	MaxOpenConn int    `json:"maxOpenConn"`
	MaxLifeTime int64  `json:"maxLifeTime"`
}

func Initialization(config *Config) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?%v",
		config.Username,
		config.Password,
		config.Address,
		config.Database,
		config.Config,
	)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         255,
		SkipInitializeWithVersion: false,
	}
	var _log zapgorm2.Logger
	if log.Instance != nil {
		_log = zapgorm2.Logger{
			ZapLogger:                 log.Instance,
			LogLevel:                  logger.Info,
			SlowThreshold:             500 * time.Millisecond,
			SkipCallerLookup:          true,
			IgnoreRecordNotFoundError: true,
		}
		_log.SetAsDefault()
	} else {
		_log.LogLevel = -1
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger:                 _log,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		inner.Fatal(err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		inner.Fatal(err.Error())
	}
	if err := sqlDB.Ping(); err != nil {
		inner.Fatal(err.Error())
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.MaxLifeTime) * time.Second)
	Instance = db
	inner.Debug("database initialization complete")
}
