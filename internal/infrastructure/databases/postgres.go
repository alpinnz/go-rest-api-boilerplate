package databases

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type PostgresConfig struct {
	Host  string
	Port  int
	User  string
	Pass  string
	Name  string
	Debug bool
}

func NewPostgres(p PostgresConfig) (*gorm.DB, error) {
	level := logger.Silent
	if p.Debug == true {
		level = logger.Info
	}
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  level, // using level from parameter
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		p.Host, p.User, p.Pass, p.Name, p.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
		Logger:         gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
