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

func NewDatabase(v *viper.Viper, log *logrus.Logger, isTest bool) *gorm.DB {
	// select database
	var dbConfig *viper.Viper
	if isTest {
		dbConfig = v.Sub("testing.database")
		if dbConfig == nil {
			log.Fatalf("testing.database configuration not found")
		}
	} else {
		dbConfig = v.Sub("database")
		if dbConfig == nil {
			log.Fatalf("database configuration not found")
		}
	}

	// Get configuration database
	username := dbConfig.GetString("username")
	password := dbConfig.GetString("password")
	host := dbConfig.GetString("host")
	port := dbConfig.GetInt("port")
	database := dbConfig.GetString("name")
	idleConnection := dbConfig.GetInt("pool.idle")
	maxConnection := dbConfig.GetInt("pool.max")
	maxLifeTimeConnection := dbConfig.GetInt("pool.lifetime")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, username, password, database, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
