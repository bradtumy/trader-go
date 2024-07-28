package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func initializeDatabase(cfg Config) error {
	var err error
	dbConnectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Database.Username, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.Name)

	db, err = sql.Open("mysql", dbConnectString)
	if err != nil {
		return err
	}
	fmt.Println("MySQL DB Connection Established")
	return nil
}
