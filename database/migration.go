package database

import (
	"cyberpull.com/go-cyb/dbo"
)

var migrations = make([]interface{}, 0)

func Migrate(db *dbo.TxDB, seed ...bool) (err error) {
	for _, model := range migrations {
		err = db.AutoMigrate(model)

		if err != nil {
			return
		}
	}

	if len(seed) > 0 && seed[0] {
		err = Seed(db)
	}

	return
}

func Migrations(models ...interface{}) {
	migrations = append(migrations, models...)
}
