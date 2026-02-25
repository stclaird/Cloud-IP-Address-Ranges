package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func IPv4ToInt(ip string) int64 {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0
	}
	ip4 := parsedIP.To4()
	if ip4 == nil {
		return 0
	}
	return int64(ip4[0])*16777216 + int64(ip4[1])*65536 + int64(ip4[2])*256 + int64(ip4[3])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_int_query.go <IP_ADDRESS>")
		fmt.Println("Example: go run test_int_query.go 13.43.34.206")
		os.Exit(1)
	}

	ipAddress := os.Args[1]
	ipInt := IPv4ToInt(ipAddress)

	fmt.Printf("🔍 Querying for IP: %s\n", ipAddress)
	fmt.Printf("   IP as integer: %d\n\n", ipInt)

	// Open database
	db, err := sql.Open("sqlite3", "cmd/main/db-output/cloudIP.sqlite3.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check schema
	var schema string
	err = db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='net'").Scan(&schema)
	if err != nil {
		log.Fatal("Failed to get schema:", err)
	}
	fmt.Println("📋 Database Schema:")
	fmt.Println(schema)
	fmt.Println()

	// Query for the IP
	query := `
		SELECT cloudplatform, net, start_ip, end_ip, iptype
		FROM net
		WHERE start_ip <= ?
		AND end_ip >= ?
		LIMIT 10
	`

	fmt.Println("🔎 SQL Query:")
	fmt.Printf("   SELECT cloudplatform, net, start_ip, end_ip, iptype\n")
	fmt.Printf("   FROM net\n")
	fmt.Printf("   WHERE start_ip <= %d\n", ipInt)
	fmt.Printf("   AND end_ip >= %d\n\n", ipInt)

	rows, err := db.Query(query, ipInt, ipInt)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var cloudplatform, netCidr, iptype string
		var startIP, endIP int

		err = rows.Scan(&cloudplatform, &netCidr, &startIP, &endIP, &iptype)
		if err != nil {
			log.Fatal("Scan failed:", err)
		}

		count++
		fmt.Printf("✅ Match #%d:\n", count)
		fmt.Printf("   Cloud Platform: %s\n", cloudplatform)
		fmt.Printf("   CIDR Range:     %s\n", netCidr)
		fmt.Printf("   Start IP:       %d\n", startIP)
		fmt.Printf("   End IP:         %d\n", endIP)
		fmt.Printf("   IP Type:        %s\n\n", iptype)
	}

	if count == 0 {
		fmt.Println("❌ No matching cloud provider ranges found")
	} else {
		fmt.Printf("📊 Total matches: %d\n", count)
	}

	// Show record count
	var totalCount, ipv4Count int
	db.QueryRow("SELECT COUNT(*), SUM(CASE WHEN iptype = 'IPv4' THEN 1 ELSE 0 END) FROM net").Scan(&totalCount, &ipv4Count)
	fmt.Printf("\n📈 Database Stats: %s total records (all IPv4)\n", strconv.Itoa(totalCount))
}
