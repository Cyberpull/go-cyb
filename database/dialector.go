package database

import (
	"fmt"
	"os"

	"cyberpull.com/go-cyb/dbo"
	"cyberpull.com/go-cyb/errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dialector() (conn gorm.Dialector, err error) {
	driver := dbo.Driver()

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	switch driver {
	case dbo.DRIVER_MYSQL:
		conn = mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUsername,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		))

	case dbo.DRIVER_PGSQL:
		conn = postgres.Open(fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			dbUsername,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		))

	default:
		err = errors.New("DB Driver not available")
	}

	return
}
