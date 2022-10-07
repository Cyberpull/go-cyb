package database

import "cyberpull.com/go-cyb/dbo"

type SeederHandler func(db *dbo.TxDB) error

var seederHandlers = make([]SeederHandler, 0)

func Seed(db *dbo.TxDB) (err error) {
	for _, handler := range seederHandlers {
		tx := db.NewSession()

		if err = handler(tx); err != nil {
			return
		}
	}

	return
}

func Seeders(handlers ...SeederHandler) {
	seederHandlers = append(seederHandlers, handlers...)
}
