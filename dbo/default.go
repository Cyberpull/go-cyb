package dbo

import (
	"cyberpull.com/go-cyb/log"

	"gorm.io/gorm"
)

type TxScope func(*gorm.DB) (tx *gorm.DB)
type TxFunction func(tx *TxDB) error

type TxDB struct {
	*gorm.DB

	opts *Options
}

func (tx *TxDB) New(v *gorm.DB) *TxDB {
	return New(v, tx.opts)
}

func (tx *TxDB) Scopes(funcs ...TxScope) *TxDB {
	scopes := make([]func(*gorm.DB) (tx *gorm.DB), 0)

	for _, scope := range funcs {
		scopes = append(scopes, scope)
	}

	tx2 := tx.DB.Scopes(scopes...)
	return tx.New(tx2)
}

func (tx *TxDB) Transaction(fn TxFunction) error {
	return tx.DB.Transaction(func(tx1 *gorm.DB) (err error) {
		newTx := New(tx1, tx.opts)

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

func New(tx *gorm.DB, opts *Options) *TxDB {
	return &TxDB{
		DB:   tx,
		opts: opts,
	}
}
