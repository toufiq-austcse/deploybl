package mysql

import (
	"database/sql"
	"time"

	"github.com/toufiq-austcse/deployit/config"
)

func getDataSourceName() string {
	// return "root:pass@tcp(localhost:3306)/entimport?parseTime=True"
	return config.AppConfig.DB_CONFIG.USER + ":" +
		config.AppConfig.DB_CONFIG.PASSWORD +
		"@tcp(" + config.AppConfig.DB_CONFIG.HOST +
		":" + config.AppConfig.DB_CONFIG.PORT + ")/" +
		config.AppConfig.DB_CONFIG.DB_NAME + "?parseTime=True"
}

func OpenMySqlDbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", getDataSourceName())
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
