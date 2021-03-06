package database

import (
	"cyberpull.com/go-cyb/dbo"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Connect(opts dbo.Options) (value *dbo.TxDB, err error) {
	var (
		db   *gorm.DB
		conn gorm.Dialector
	)

	var config *gorm.Config = &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	conn, err = dialector(&opts)

	if err != nil {
		return
	}

	db, err = gorm.Open(conn, config)

	if err != nil {
		return
	}

	if dbo.Driver(&opts) == dbo.DRIVER_PGSQL {
		err = db.Exec(`SET DEFAULT_TRANSACTION_ISOLATION TO SERIALIZABLE`).Error
		// SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
		// SET DEFAULT_TRANSACTION_ISOLATION TO SERIALIZABLE;

		if err != nil {
			return
		}
	}

	value = dbo.New(db, &opts)

	return
}
