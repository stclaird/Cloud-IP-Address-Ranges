package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	testIP := "8.8.8.8"
	
	fmt.Printf("=== Querying for IP: %s ===\n\n", testIP)

	// Test SQLite
	fmt.Println("📦 SQLite Query:")
	querySQLite(testIP)

	fmt.Println()

	// Test DuckDB
	fmt.Println("🦆 DuckDB Query:")
	queryDuckDB(testIP)
}

func querySQLite(ip string) {
	db, err := sql.Open("sqlite3", "../../cmd/main/db-output/cloudIP.sqlite3.db")
	if err != nil {
		log.Printf("Error opening SQLite: %v", err)
		return
	}
	defer db.Close()

	// Text-based comparison
	query := `
		SELECT cloudplatform, net, iptype
		FROM net
		WHERE start_ip <= ?
		AND end_ip >= ?
		LIMIT 5
	`

	rows, err := db.Query(query, ip, ip)
	if err != nil {
		log.Printf("Error querying: %v", err)
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var platform, net, iptype string
		if err := rows.Scan(&platform, &net, &iptype); err != nil {
			log.Printf("Error scanning: %v", err)
			continue
		}
		fmt.Printf("  ✓ %s - %s (%s)\n", platform, net, iptype)
		found = true
	}

	if !found {
		fmt.Println("  No matches found")
	}
}

func queryDuckDB(ip string) {
	db, err := sql.Open("duckdb", "../../cmd/main/db-output/cloudIP.duckdb")
	if err != nil {
		log.Printf("Error opening DuckDB: %v", err)
		return
	}
	defer db.Close()

	// Install inet extension
	if _, err := db.Exec("INSTALL inet; LOAD inet;"); err != nil {
		log.Printf("Warning: inet extension not available: %v", err)
		// Fall back to text comparison
		queryDuckDBTextMode(db, ip)
		return
	}

	// Native IP comparison
	query := `
		SELECT cloudplatform, net, iptype
		FROM net
		WHERE TRY_CAST(start_ip AS INET) <= TRY_CAST(? AS INET)
		AND TRY_CAST(end_ip AS INET) >= TRY_CAST(? AS INET)
		LIMIT 5
	`

	rows, err := db.Query(query, ip, ip)
	if err != nil {
		log.Printf("Error querying: %v", err)
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var platform, net, iptype string
		if err := rows.Scan(&platform, &net, &iptype); err != nil {
			log.Printf("Error scanning: %v", err)
			continue
		}
		fmt.Printf("  ✓ %s - %s (%s)\n", platform, net, iptype)
		found = true
	}

	if !found {
		fmt.Println("  No matches found")
	}
}

func queryDuckDBTextMode(db *sql.DB, ip string) {
	query := `
		SELECT cloudplatform, net, iptype
		FROM net
		WHERE start_ip <= ?
		AND end_ip >= ?
		LIMIT 5
	`

	rows, err := db.Query(query, ip, ip)
	if err != nil {
		log.Printf("Error querying: %v", err)
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var platform, net, iptype string
		if err := rows.Scan(&platform, &net, &iptype); err != nil {
			log.Printf("Error scanning: %v", err)
			continue
		}
		fmt.Printf("  ✓ %s - %s (%s) [text mode]\n", platform, net, iptype)
		found = true
	}

	if !found {
		fmt.Println("  No matches found")
	}
}
