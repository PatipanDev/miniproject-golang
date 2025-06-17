package configs

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DB_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	log.Println("connect the database >>>>>")

	return db, nil
}

func NewDatabaseGromRiver() (*gorm.DB, error) {
	sqlDB, err := sql.Open("pgx", DB_URL)
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
