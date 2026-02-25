package main

import (
	"context"
	"database/sql"
	"fmt"
	"flag"
	"log"
	"os"

	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/config"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/ipfile"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/ipnet"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/model"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/repository"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/marcboeker/go-duckdb"
)

var confObj config.Config
var databases map[string]*sql.DB

func initDatabase(dbType string) (*sql.DB, error) {
	var dbPath string
	var driver string
	
	switch dbType {
	case "sqlite":
		driver = "sqlite3"
		dbPath = fmt.Sprintf("%s/%s.sqlite3.db", confObj.Dbdir, confObj.Dbfile)
		// SQLite requires file to be created first
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("os Create Error: %w", err)
		}
		file.Close()
	case "duckdb":
		driver = "duckdb"
		dbPath = fmt.Sprintf("%s/%s.duckdb", confObj.Dbdir, confObj.Dbfile)
		// DuckDB creates its own file, so remove if exists first
		os.Remove(dbPath)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Open database connection
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s database: %w", dbType, err)
	}

	return db, nil
}

func init() {
	confObj = config.NewConfig()

	if confObj.Downloaddir == "" {
		confObj.Downloaddir = "downloadedfiles"
	}

	err := os.MkdirAll(confObj.Downloaddir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(confObj.Dbdir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize databases map
	databases = make(map[string]*sql.DB)

	// Initialize each requested database type
	for _, dbType := range confObj.Dbtypes {
		db, err := initDatabase(dbType)
		if err != nil {
			log.Fatalf("Failed to initialize %s database: %v", dbType, err)
		}
		databases[dbType] = db
		log.Printf("Initialized %s database", dbType)
	}
}

func main() {
	debugDownload := flag.Bool("debug", false, "enable debug logging for IP file downloads")
	flag.Parse()

	ctx := context.Background()

	// Close all databases on exit
	defer func() {
		for dbType, db := range databases {
			log.Printf("Closing %s database", dbType)
			db.Close()
		}
	}()

	// Setup all databases
	for dbType, db := range databases {
		log.Printf("Setting up %s database schema", dbType)
		model.SetupDB(db)
	}

	// Create repositories for each database
	cidrRepos := make(map[string]repository.CidrRepository)
	for dbType, db := range databases {
		cidrRepos[dbType] = repository.NewCidrRepository(db)
	}

	var report []reportEntry //create a report struct to keep track of inserts

	for _, i := range confObj.Ipfiles {

		var cidrs []string

		var entry = reportEntry{ //init report struct entry for each cloud provider
			CloudPlatform: i.Cloudplatform,
			Success:       0,
			Failed:        0,
		}

		downloadto := fmt.Sprintf("%s/%s", confObj.Downloaddir, i.DownloadFileName)

		var url string
		url = i.Url

		var err error
		switch i.Cloudplatform {
		case "azure":
			url, err = ipfile.ResolveAzureDownloadUrl(*debugDownload) //azure download file changes so we need to figure out what the latest path is
		}

		if err != nil {
			break
		}

		FileObj := ipfile.IpfileTXT{
			Common: ipfile.Common{
				Debug: *debugDownload,
			},
		}
		FileObj.Download(downloadto, url)
		cidrs_raw := ipfile.AsText[ipfile.IpfileTXT](downloadto)
		cidrs = ipfile.Process(cidrs_raw)

		for _, cidr := range cidrs {
			processedCidr, err := ipnet.PrepareCidrforDB(cidr)
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}

			// Only IPv4 addresses
			c := model.Cidr{
				Net:           cidr,
				Start_ip:      processedCidr.NetIPDecimal,
				End_ip:        processedCidr.BcastIPDecimal,
				Url:           i.Url,
				Cloudplatform: i.Cloudplatform,
				Iptype:        processedCidr.Iptype,
			}
			
			// Insert into all configured databases
			for dbType, cidrRepo := range cidrRepos {
				_, exists := cidrRepo.FindByNet(ctx, c.Net)
				if !exists {
					err := cidrRepo.Insert(ctx, c)
					if err != nil {
						if dbType == confObj.Dbtypes[0] { // Only count once
							entry.IncrementFailed()
						}
						log.Printf("Failed to insert into %s: %v", dbType, err)
					} else {
						if dbType == confObj.Dbtypes[0] { // Only count once
							entry.IncrementSuccess()
						}
					}
				}
			}
		}
		report = append(report, entry)
		ipfile.WriteFile(downloadto, cidrs) //overwrite downloaded file with IP address info only
	}

	printReport(report)
}
