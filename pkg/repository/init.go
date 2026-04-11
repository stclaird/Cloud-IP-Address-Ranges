package repository

import (
	"context"
	"database/sql"

	model "github.com/stclaird/Cloud-IP-Address-Ranges/pkg/model"
)

type CidrRepo struct {
	DB     *sql.DB
	DbType string
}

func NewCidrRepository(db *sql.DB, dbType string) CidrRepository {
	return &CidrRepo{
		DB:     db,
		DbType: dbType,
	}
}

type CidrRepository interface {
	Insert(ctx context.Context, cidr model.Cidr) error
	FindByNet(ctx context.Context, net string) (error, bool)
}
