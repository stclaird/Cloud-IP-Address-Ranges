package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type CidrObject struct {
	Net           string
	Start_ip      int
	End_ip        int
	Url           string
	Cloudplatform string
	Iptype        string
}

const createNetDB string = `
CREATE TABLE IF NOT EXISTS net (
	net_id INTEGER PRIMARY KEY,
	net TEXT NOT NULL,
	start_ip INT NOT NULL,
	end_ip INT NOT NULL,
	url TEXT NOT NULL,
	cloudplatform TEXT NOT NULL,
	iptype TEXT NOT NULL
	);`


func SetupDB(db *sql.DB){
	_, err := db.Exec(createNetDB)
	if err != nil {
		log.Fatal(err)
	}
}

func AddCidr(db *sql.DB, newCidr CidrObject) {
	//Add a Cidr
	stmt, err := db.Prepare("INSERT INTO net ( net_id, net, start_ip, end_ip, url, cloudplatform, iptype) VALUES ( ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Prepare Error: %s %s %s", newCidr.Cloudplatform, newCidr.Net, err)
		return
	}
	res, err := stmt.Exec(newCidr.Start_ip, newCidr.Net, newCidr.Start_ip, newCidr.End_ip, newCidr.Url, newCidr.Cloudplatform, newCidr.Iptype)
	if err != nil {
		log.Printf("Insert Error: %s %s %s", newCidr.Cloudplatform, newCidr.Net, err)
		return
	}

	res.LastInsertId()
	return
}

func GetCidr(db *sql.DB) {
	//Fetch a row
	rows, err := db.Query("SELECT net, start_ip, end_ip, url, cloudplatform, iptype from net;")

	for rows.Next() {
		res := CidrObject{}
		err = rows.Scan(&res.Net, &res.Start_ip, &res.End_ip, &res.Url, &res.Cloudplatform, &res.Iptype)

		if err != nil {
			log.Fatal(err)
		}

		rows.Close()

	}
}
