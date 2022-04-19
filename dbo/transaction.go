package dbo

import (
	"cyberpull.com/go-cyb/log"

	"gorm.io/gorm"
)

type TxFunction func(tx *TxDB) error

type TxDB struct {
	*gorm.DB
}

func (tx *TxDB) Transaction(fn TxFunction) error {
	return tx.DB.Transaction(func(tx1 *gorm.DB) (err error) {
		newTx := New(tx1)

		defer func() {
			if err != nil {
				log.Errorfln("Tx Error: %s", err)
			}
		}()

		// if driver := Driver(); driver == DRIVER_PGSQL {
		// 	err = newTx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`).Error

		// 	if err != nil {
		// 		return
		// 	}

		// 	err = newTx.Exec(`SET CONSTRAINTS ALL DEFERRED`).Error

		// 	if err != nil {
		// 		return
		// 	}
		// }

		err = fn(newTx)

		return
	})
}

/********************************/

func New(tx *gorm.DB) *TxDB {
	return &TxDB{
		DB: tx,
	}
}
