package main

import (
	"context"
	"fmt"
	"flag"
	"log"
	"os"

	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/config"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/database"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/ipfile"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/ipnet"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/model"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/repository"
)

var confObj config.Config

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
}

func main() {
	debugDownload := flag.Bool("debug", false, "enable debug logging for IP file downloads")
	flag.Parse()

	ctx := context.Background()

	// Download and process IP data once
	var allCidrs []model.Cidr
	var report []reportEntry

	log.Println("Downloading and processing IP data...")
	
	for _, i := range confObj.Ipfiles {
		var cidrs []string

		var entry = reportEntry{
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
			url, err = ipfile.ResolveAzureDownloadUrl(*debugDownload)
		}

		if err != nil {
			log.Printf("Error resolving URL for %s: %v", i.Cloudplatform, err)
			continue
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
				entry.IncrementFailed()
				continue
			}

			c := model.Cidr{
				Net:           cidr,
				Start_ip:      processedCidr.NetIPString,
				End_ip:        processedCidr.BcastIPString,
				Url:           i.Url,
				Cloudplatform: i.Cloudplatform,
				Iptype:        processedCidr.Iptype,
			}
			allCidrs = append(allCidrs, c)
			entry.IncrementSuccess()
		}
		
		report = append(report, entry)
		ipfile.WriteFile(downloadto, cidrs)
	}

	log.Printf("Processed %d CIDR records", len(allCidrs))

	// Now insert into each configured database
	for _, dbConfig := range confObj.Databases {
		log.Printf("\n=== Creating %s database: %s ===", dbConfig.Type, dbConfig.Filename)
		
		fullPath := fmt.Sprintf("%s%s", confObj.Dbdir, dbConfig.Filename)
		
		// Create database file
		file, err := os.Create(fullPath)
		if err != nil {
			log.Printf("Error creating %s: %v", fullPath, err)
			continue
		}
		file.Close()

		// Connect to database
		db, err := database.NewDatabase(database.DBConfig{
			Type:     database.DBType(dbConfig.Type),
			FilePath: fullPath,
		})
		if err != nil {
			log.Printf("Error opening %s database: %v", dbConfig.Type, err)
			continue
		}

		// Setup schema
		if err := db.Setup(); err != nil {
			log.Printf("Error setting up %s database: %v", dbConfig.Type, err)
			db.Close()
			continue
		}

		// Insert data
		cidrRepo := repository.NewCidrRepository(db.GetDB())
		
		inserted := 0
		failed := 0
		
		for _, c := range allCidrs {
			_, exists := cidrRepo.FindByNet(ctx, c.Net)
			if !exists {
				if err := cidrRepo.Insert(ctx, c); err != nil {
					failed++
				} else {
					inserted++
				}
			}
		}

		db.Close()
		
		log.Printf("✓ %s: Inserted %d records, Failed %d", dbConfig.Filename, inserted, failed)
	}

	log.Println("\n=== Summary Report ===")
	printReport(report)
}
