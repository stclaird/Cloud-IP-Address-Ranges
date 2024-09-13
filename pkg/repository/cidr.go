package repository

import (
	"context"
	"log"

	model "github.com/stclaird/Cloud-IP-Address-Ranges/pkg/model"
)

func (repo *CidrRepo) Insert(ctx context.Context, newCidr model.Cidr) error {
	stmt, err := repo.DB.Prepare("INSERT INTO net ( net, start_ip, end_ip, url, cloudplatform, iptype) VALUES ( ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Prepare Error: %s %s %s", newCidr.Cloudplatform, newCidr.Net, err)
		return err
	}
	res, err := stmt.Exec( newCidr.Net, newCidr.Start_ip, newCidr.End_ip, newCidr.Url, newCidr.Cloudplatform, newCidr.Iptype)
	if err != nil {
		log.Printf("Insert Error: %s %s %s", newCidr.Cloudplatform, newCidr.Net, err)
	return err
	}

	res.LastInsertId()
	return err
}
