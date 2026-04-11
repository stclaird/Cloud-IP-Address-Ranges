package repository

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	model "github.com/stclaird/Cloud-IP-Address-Ranges/pkg/model"
)

func (repo *CidrRepo) Insert(ctx context.Context, newCidr model.Cidr) error {
	// Prepare appropriate INSERT depending on DB type. For DuckDB use CAST(... AS inet)
	var stmt *sql.Stmt
	var err error
	if repo.DbType == "duckdb" {
		stmt, err = repo.DB.Prepare("INSERT INTO net ( net, start_ip, end_ip, url, cloudplatform, iptype) VALUES ( ?, CAST(? AS inet), CAST(? AS inet), ?, ?, ?)")
	} else {
		stmt, err = repo.DB.Prepare("INSERT INTO net ( net, start_ip, end_ip, url, cloudplatform, iptype) VALUES ( ?, ?, ?, ?, ?, ?)")
	}
	if err != nil {
		log.Printf("Prepare Error: %s %s %s", newCidr.Cloudplatform, newCidr.Net, err)
		return err
	}

	if repo.DbType == "duckdb" {
		// For DuckDB with inet type insert textual IP forms
		start := newCidr.Start_ip_text
		end := newCidr.End_ip_text
		if start == "" {
			// fallback to dotted/colon notation from numeric if available
			start = strconv.FormatInt(newCidr.Start_ip, 10)
		}
		if end == "" {
			end = strconv.FormatInt(newCidr.End_ip, 10)
		}
		res, err := stmt.Exec(newCidr.Net, start, end, newCidr.Url, newCidr.Cloudplatform, newCidr.Iptype)
		if err != nil {
			log.Printf("Insert Error: %s %s %s", newCidr.Cloudplatform, newCidr.Net, err)
			return err
		}
		_, _ = res.LastInsertId()
		return nil
	}

	res, err := stmt.Exec(newCidr.Net, newCidr.Start_ip, newCidr.End_ip, newCidr.Url, newCidr.Cloudplatform, newCidr.Iptype)
	if err != nil {
		log.Printf("Insert Error: %s %s %s", newCidr.Cloudplatform, newCidr.Net, err)
		return err
	}

	_, _ = res.LastInsertId()
	return nil
}

func (repo *CidrRepo) FindByNet(ctx context.Context, net string) (error, bool) {
	stmt, err := repo.DB.Prepare("SELECT net FROM net where net=?")
	if err != nil {
		log.Printf("Prepare Error")
		return err, false
	}

	var fetchedCidr model.Cidr

	err = stmt.QueryRow(net).Scan(&fetchedCidr.Net)
	if err != nil {
		if err == sql.ErrNoRows {
			//Row not found
			return err, false
		}
	}
	return err, true
}
