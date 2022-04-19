package dbo

import "os"

const (
	DRIVER_MYSQL string = "mysql"
	DRIVER_PGSQL string = "pgsql"
)

func Driver() string {
	driver := os.Getenv("DB_DRIVER")

	if driver == "" {
		driver = DRIVER_PGSQL
	}

	return driver
}
