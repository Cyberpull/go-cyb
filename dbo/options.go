package dbo

type DRIVER string

const (
	DRIVER_MYSQL DRIVER = "mysql"
	DRIVER_PGSQL DRIVER = "pgsql"
)

type Options struct {
	Driver   DRIVER
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	URL      string
}
