package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"cyberpull.com/go-cyb/database"
	"cyberpull.com/go-cyb/dbo"
	"cyberpull.com/go-cyb/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ExampleModel struct {
	gorm.Model

	Name string `gorm:"index;column:name;"`
}

type DBOTestSuite struct {
	suite.Suite

	db    *dbo.TxDB
	total uint
}

func (s *DBOTestSuite) SetupSuite() {
	s.total = 300

	var err error

	database.Migrations(&ExampleModel{})
	database.Seeders(s.seeder())

	options := dbo.Options{
		Driver:   dbo.DRIVER(os.Getenv("DB_DRIVER")),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DATABASE"),
	}

	s.db, err = database.Connect(options)
	require.NoError(s.T(), err)

	err = database.Migrate(s.db, true)
	require.NoError(s.T(), err)

	log.Println("Database Migration Successful!")
}

func (s *DBOTestSuite) seeder() database.SeederHandler {
	return func(db *dbo.TxDB) (err error) {
		stmt := gorm.Statement{DB: db.DB}
		stmt.Parse(&ExampleModel{})

		table := stmt.Schema.Table
		query := fmt.Sprintf("TRUNCATE TABLE %s", table)

		if err = s.db.Exec(query).Error; err != nil {
			return
		}

		entries := make([]*ExampleModel, 0)

		for i := 0; i < int(s.total); i++ {
			entry := &ExampleModel{}
			entry.ID = uint(i + 1)
			entry.Name = fmt.Sprintf("Model %d", entry.ID)

			entries = append(entries, entry)
		}

		err = db.Create(&entries).Error

		return
	}
}

func (s *DBOTestSuite) TestPaginate() {
	limit := 10

	tx := s.db.Order(`"id" DESC`)

	data, err := dbo.Paginate[*ExampleModel](tx, 1, uint(limit))
	require.NoError(s.T(), err)

	// Assert count
	assert.Len(s.T(), data.Data, limit)

	// Assert current page
	assert.Equal(s.T(), uint(1), data.Current_page)

	// Assert last page
	assert.Equal(s.T(), uint(s.total/uint(limit)), data.Last_page)

	// Assert total records
	assert.Equal(s.T(), s.total, data.Total)
}

func (s *DBOTestSuite) TestNullString() {
	var data dbo.Null[string]

	data.Scan("Hello")

	assert.Equal(s.T(), data.Data, "Hello")
	assert.True(s.T(), data.Valid)

	data.Scan("")

	assert.Equal(s.T(), data.Data, "")
	assert.False(s.T(), data.Valid)
}

func (s *DBOTestSuite) TestNullUint() {
	var data dbo.Null[uint]

	data.Scan(200)

	assert.Equal(s.T(), data.Data, uint(200))
	assert.True(s.T(), data.Valid)

	data.Scan(0)

	assert.Equal(s.T(), data.Data, uint(0))
	assert.False(s.T(), data.Valid)
}

func (s *DBOTestSuite) TestNullTime() {
	var data dbo.Null[time.Time]

	now := time.Now()
	data.Scan(now)

	assert.Equal(s.T(), data.Data, now)
	assert.True(s.T(), data.Valid)

	now = time.Time{}
	data.Scan(now)

	assert.Equal(s.T(), data.Data, now)
	assert.False(s.T(), data.Valid)
}

/******************************************/

func TestDBO(t *testing.T) {
	suite.Run(t, new(DBOTestSuite))
}
