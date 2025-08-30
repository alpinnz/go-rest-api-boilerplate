package databases

import (
	"fmt"

	customLogger "github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
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
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		p.Host, p.User, p.Pass, p.Name, p.Port,
	)

	var gormLogger logger.Interface
	if p.Debug {
		gormLogger = customLogger.NewGormLogger(logger.Info)
	} else {
		gormLogger = customLogger.NewGormLogger(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
		Logger:         gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
