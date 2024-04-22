package main

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketName = "functions"
	location   = "eu-west-1"
)

type StorageService struct {
	Client *minio.Client
}

func StartStorageConnection(endpoint string, accessKeyID string, secretAccessKey string, storageSecure string) (storageService StorageService, err error) {
	client, err := minio.New(

		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: storageSecure == "true",
		},
	)

	storageService = StorageService{Client: client}

	return
}

func (s *StorageService) UploadFile(objectName string, filePath string) (err error) {
	_, err = s.Client.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	return
}
