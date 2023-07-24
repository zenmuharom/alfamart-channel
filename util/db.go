package util

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func ConnectDB() error {

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DB_User, config.DB_Pass, config.DB_Address, config.DB_Port, config.DB_Name)
	client, err := sqlx.Connect("mysql", dataSource)
	if err != nil {
		return errors.New(fmt.Sprintf("error: %s source: %s", err.Error(), dataSource))
	}
	client.SetConnMaxLifetime(time.Minute)
	client.SetMaxIdleConns(40)
	client.SetMaxOpenConns(100)
	if err != nil {
		return errors.New(fmt.Sprintf("error: %s source: %s", err.Error(), dataSource))
		return err
	}

	db = client

	return nil
}

func GetDB() *sqlx.DB {
	return db
}
