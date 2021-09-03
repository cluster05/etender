package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "root"
	hostname = "127.0.0.1:3306"
	dbName   = "tenderDB"
)

func MysqlDB() *sql.DB {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName))
	if err != nil {
		log.Printf("[MYSQL] Error %s when opening DB\n", err)
		return nil
	}

	/*
		do not call here otherwise it will close
		the movement to called this function and retrun db
	*/
	// defer db.Close()

	return db
}
