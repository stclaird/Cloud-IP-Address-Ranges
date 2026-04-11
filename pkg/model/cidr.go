package model

import (
	"database/sql"
	"log"
)

type Cidr struct {
	Net           string
	Start_ip      int64
	End_ip        int64
	Start_ip_text string
	End_ip_text   string
	Url           string
	Cloudplatform string
	Iptype        string
}

const createNetDB string = `
CREATE TABLE IF NOT EXISTS net (
	net TEXT PRIMARY KEY,
	start_ip BIGINT NOT NULL,
	end_ip BIGINT NOT NULL,
	url TEXT NOT NULL,
	cloudplatform TEXT NOT NULL,
	iptype TEXT NOT NULL
	);`

const createNetDBDuckDB string = `
CREATE TABLE IF NOT EXISTS net (
	net TEXT PRIMARY KEY,
	start_ip TEXT NOT NULL,
	end_ip TEXT NOT NULL,
	url TEXT NOT NULL,
	cloudplatform TEXT NOT NULL,
	iptype TEXT NOT NULL
	);`

func SetupDB(db *sql.DB, dbType string) {
	var err error
	if dbType == "duckdb" {
		_, err = db.Exec(createNetDBDuckDB)
	} else {
		_, err = db.Exec(createNetDB)
	}
	if err != nil {
		log.Fatal(err)
	}
}
