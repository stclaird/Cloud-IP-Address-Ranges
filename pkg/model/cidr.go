package model

import (
	"database/sql"
	"log"
)

type Cidr struct {
    Net string
	Start_ip      int
	End_ip        int
	Url           string
	Cloudplatform string
	Iptype        string
}

const createNetDB string = `
CREATE TABLE IF NOT EXISTS net (
	net TEXT PRIMARY KEY,
	start_ip INT NOT NULL,
	end_ip INT NOT NULL,
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
