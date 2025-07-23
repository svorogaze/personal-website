package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

const (
	minPasswordLength    = 5
	maxPasswordLength    = 25
	minLoginLength       = 5
	maxLoginLength       = 25
	fileSizeLimit        = 5 << 20 // 5 MB
	minTitleLength       = 1
	maxTitleLength       = 75
	minDescriptionLength = 1
	maxDescriptionLength = 75
)

type API struct {
	db          *sqlx.DB
	minioClient *minio.Client
}

func New(db *sqlx.DB, minioClient *minio.Client) *API {
	return &API{db: db, minioClient: minioClient}
}
