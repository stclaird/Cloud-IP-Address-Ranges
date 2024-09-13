package repository

import (
	"context"
	"database/sql"

	model "github.com/stclaird/Cloud-IP-Address-Ranges/pkg/model"
)

type CidrRepo struct {
	DB *sql.DB
}

func NewCidrRepository(db *sql.DB) CidrRepository {
	return &CidrRepo{
		DB: db,
	}
}

type CidrRepository interface {
	Insert(ctx context.Context, cidr model.Cidr) error
}
