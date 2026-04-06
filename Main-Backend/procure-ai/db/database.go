package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ConnectionString = "postgresql://postgres:LowKey7642@localhost:5432/procure_ai"

const dsn = "host=localhost user=postgres password=LowKey7642 dbname=procure_ai port=5432 sslmode=disable"

func NewPostgresDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
