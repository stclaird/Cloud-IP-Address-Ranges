package model

import (
	"database/sql"
	"log"
)

type Cidr struct {
    Net string
	Start_ip      string
	End_ip        string
	Url           string
	Cloudplatform string
	Iptype        string
}

const createNetDB string = `
CREATE TABLE IF NOT EXISTS net (
	net TEXT PRIMARY KEY,
	start_ip TEXT NOT NULL,
	end_ip TEXT NOT NULL,
	url TEXT NOT NULL,
	cloudplatform TEXT NOT NULL,
	iptype TEXT NOT NULL
	);`


func SetupDB(db *sql.DB) {
    _, err := db.Exec(createNetDB)
    if err != nil {
        log.Fatal(err)
    }
}
