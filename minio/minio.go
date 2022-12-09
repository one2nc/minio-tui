package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint, AccessKey, SecretAccessKey string
	UseSSL                               bool
}

func GetMinioClient(config *Config) (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}

func GetBuckets(client *minio.Client) ([]minio.BucketInfo, error) {
	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		fmt.Println(err)
		return make([]minio.BucketInfo, 0), err
	}
	return buckets, nil
}

func GetFiles(bucketName string, client *minio.Client) ([]minio.ObjectInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})

	objects := []minio.ObjectInfo{}
	for object := range objectCh {
		if object.Err != nil {
			return objects, object.Err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func MakeBucket(bucketName string,client *minio.Client,bucketOptions minio.MakeBucketOptions) error {
	err := client.MakeBucket(context.Background(),bucketName,bucketOptions)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Bucket Created Succesfully!!!!") 
	return nil
}
