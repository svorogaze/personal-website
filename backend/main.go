package main

import (
	"backend/api"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func createBucket(minioClient *minio.Client, bucketName string) error {
	ctx := context.Background()
	b, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !b {
		err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		err = minioClient.SetBucketPolicy(ctx, bucketName,
			`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Effect": "Allow",
							"Principal": {"AWS": ["*"]},
							"Action": ["s3:GetObject"],
							"Resource": ["arn:aws:s3:::`+bucketName+`/*"]
						}
					]
		}`)
		return err
	}
	return nil
}

func main() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", "db", os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close db: %v", err)
		}
	}(db)

	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to connect to minio %v", err)
	}
	err = createBucket(minioClient, "blog-cover-images")
	if err != nil {
		log.Fatalf("failed to create the image bucket %v", err)
	}

	API := api.New(db, minioClient)

	router := gin.Default()
	router.GET("/api/blogs/:id", API.GetBlog)
	router.GET("/api/blogs", API.GetBlogsRange)
	router.POST("/api/blogs", API.CreateBlog)
	err = router.Run(":8088")
	if err != nil {
		log.Fatalf("Error when running server: %v", err)
	}
}
