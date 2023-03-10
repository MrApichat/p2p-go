package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func Con() {
	var err error
	// connection := "postgres://postgres:12345678@database-1.ckhl0w0gyple.ap-northeast-1.rds.amazonaws.com/p2p_go"
	Db, err = sql.Open("postgres", "postgres://postgres:123456@localhost/p2p-go?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = Db.Ping()

	if err != nil {
		log.Fatal(err)
	}
}
