package objectstorage

import (
	"context"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var	(
  bucketName = "testbucket"
	location = "us-east-1"
  minioClient *minio.Client
)

//FileName is fileheader.name + id

func Init() error {
	endpoint := "play.min.io"
	accessKeyID := "22rojoiu8ASDHY(@*!SF"
	secretAccessKey := "zasd+=sdpASD*@sd*02enlitnifILbZam1XMASID"
	useSSL := true

  var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
  return err
}

func Upload(ctx context.Context, name string, file multipart.File, fheader multipart.FileHeader) error {
  defer file.Close()

  err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
      return err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	_, err = minioClient.PutObject(ctx, bucketName, name, file, fheader.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"} )
  return err
}

func Remove(ctx context.Context, name string) error {
  err := minioClient.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions {ForceDelete: true})
	return err
}

func Get(ctx context.Context, name string) (*url.URL, error) {
  link, err := minioClient.PresignedGetObject(ctx, bucketName, name,time.Hour, nil)
	return link, err
}
