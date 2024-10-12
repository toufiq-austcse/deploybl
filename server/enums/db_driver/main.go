package db_drivers_enum

type (
	DbDriver string
	CREATED  int64
)

const (
	POSTGRESQL DbDriver = "postgres"
	MYSQL      DbDriver = "mysql"
)
