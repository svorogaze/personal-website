package api

import (
	"context"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
)

func (api *API) uploadImage(fh *multipart.FileHeader, bucketName string, filename string) error {
	cxt := context.Background()
	f, err := fh.Open()
	if err != nil {
		return err
	}
	_, err = api.minioClient.PutObject(cxt, bucketName, filename, f, fh.Size, minio.PutObjectOptions{})
	return err
}
