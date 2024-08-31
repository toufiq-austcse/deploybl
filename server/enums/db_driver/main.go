package db_drivers_enum

type DbDriver string
type CREATED int64

const (
	POSTGRESQL DbDriver = "postgres"
	MYSQL      DbDriver = "mysql"
)
