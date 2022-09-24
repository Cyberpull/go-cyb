package dbo

import "cyberpull.com/go-cyb/errors"

type DBFunction func(tx ...*TxDB) (value *TxDB, err error)

func DB(db *TxDB) DBFunction {
	// Database Function
	return func(tx ...*TxDB) (value *TxDB, err error) {
		if len(tx) > 0 && tx[0] != nil {
			value = tx[0]
			return
		}

		if db != nil {
			value = db.NewSession()
			return
		}

		err = errors.New("Database connection not found")

		return
	}
}
