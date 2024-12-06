package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Postgres struct {
	db                 *gorm.DB
	MaxIdleConnections int
	MaxOpenConnections int
}

func InitializeDBPostgres(maxIdleConnections, maxOpenConnections int) *Postgres {
	postgresDB := Postgres{
		MaxIdleConnections: maxIdleConnections,
		MaxOpenConnections: maxOpenConnections,
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connectionDBUrl := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s`, dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(connectionDBUrl), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	postgresDB.db = db

	sqlDB, err := postgresDB.db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(postgresDB.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(postgresDB.MaxOpenConnections)

	postgresDB.db = db
	log.Info("Connected to Postgres DB")
	return &postgresDB
}

func (postgresDB *Postgres) GetDB() *gorm.DB {
	return postgresDB.db
}

func (postgresDB *Postgres) PingDB() error {
	db, err := postgresDB.GetDB().DB()
	if err != nil {
		return err
	}
	return db.Ping()
}
