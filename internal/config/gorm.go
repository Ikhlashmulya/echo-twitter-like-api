package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(config *viper.Viper, log *logrus.Logger) *gorm.DB {
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	dbname := config.GetString("database.name")
	idleConnection := config.GetInt("database.pool.idle")
	maxConnection := config.GetInt("database.pool.max")
	maxLifeTimeConnection := config.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(newLogrusWriter(log), logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("error connecting to database : %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("error connecting to database : %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func newLogrusWriter(log *logrus.Logger) *logrusWriter {
	return &logrusWriter{
		Logger: log,
	}
}

func (lw *logrusWriter) Printf(message string, ars ...any) {
	lw.Logger.Tracef(message, ars...)
}
