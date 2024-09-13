package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/config"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/ipfile"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/ipnet"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/model"
	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/repository"

	_ "github.com/mattn/go-sqlite3"
)

var confObj config.Config
var db *sql.DB

func init() {
	confObj = config.NewConfig()

	if confObj.Downloaddir == "" {
		confObj.Downloaddir = "downloadedfiles"
	}

	err := os.MkdirAll(confObj.Downloaddir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(confObj.Dbdir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	full_path := fmt.Sprintf("%s/%s", confObj.Dbdir, confObj.Dbfile)
	file, err := os.Create(full_path)

	if err != nil {
		log.Println("Os Create Error: ", err)
	}

	file.Close()

	db, err = sql.Open("sqlite3", full_path)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	ctx := context.Background()

	defer db.Close()
	model.SetupDB(db)

	cidrRepo := repository.NewCidrRepository(db)


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
			url, err = ipfile.ResolveAzureDownloadUrl() //azure download file changes so we need to figure out what the latest path is
		}

		if err != nil {
			break
		}

		var FileObj ipfile.IpfileTXT
		FileObj.Download(downloadto, url)
		cidrs_raw := ipfile.AsText[ipfile.IpfileTXT](downloadto)
		cidrs = ipfile.Process(cidrs_raw)

		for _, cidr := range cidrs {
			processedCidr, err := ipnet.PrepareCidrforDB(cidr)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			if processedCidr.Iptype == "IPv4" {
				c := model.Cidr{
					Net:           cidr,
					Start_ip:      processedCidr.NetIPDecimal,
					End_ip:        processedCidr.BcastIPDecimal,
					Url:           i.Url,
					Cloudplatform: i.Cloudplatform,
					Iptype:        processedCidr.Iptype,
				}
				_, exists := cidrRepo.FindByNet(ctx, c.Net)
				if !exists {
					err := cidrRepo.Insert(ctx,c)
					//record inserts
					if err != nil {
						entry.IncrementFailed()
					} else {
						entry.IncrementSuccess()
					}
				}
			}
		}
		report = append(report, entry)
		ipfile.WriteFile(downloadto, cidrs) //overwrite downloaded file with IP address info only
	}

	printReport(report)
}
