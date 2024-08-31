package postgresql

import (
	"database/sql"
	"github.com/toufiq-austcse/deployit/config"
	"time"
)

func getDataSourceName() string {
	return "host=" + config.AppConfig.DB_CONFIG.HOST +
		" port=" + config.AppConfig.DB_CONFIG.PORT +
		" user=" + config.AppConfig.DB_CONFIG.USER +
		" dbname=" + config.AppConfig.DB_CONFIG.DB_NAME +
		" password=" + config.AppConfig.DB_CONFIG.PASSWORD +
		" sslmode=disable"
}

func OpenPgDbConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", getDataSourceName())
	if err != nil {
		return nil, err
	}
	if pingErr := db.Ping(); pingErr != nil {
		return nil, pingErr
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return db, nil
}
