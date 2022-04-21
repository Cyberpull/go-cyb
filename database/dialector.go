package database

import (
	"fmt"

	"cyberpull.com/go-cyb/dbo"
	"cyberpull.com/go-cyb/errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dialector(opts *dbo.Options) (conn gorm.Dialector, err error) {
	switch dbo.Driver(opts) {
	case dbo.DRIVER_MYSQL:
		conn = mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			opts.Username,
			opts.Password,
			opts.Host,
			opts.Port,
			opts.DBName,
		))

	case dbo.DRIVER_PGSQL:
		conn = postgres.Open(fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			opts.Username,
			opts.Password,
			opts.Host,
			opts.Port,
			opts.DBName,
		))

	default:
		err = errors.New("DB Driver not available")
	}

	return
}
