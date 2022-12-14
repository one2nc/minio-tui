package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// path is to store objects
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

func MakeBucket(bucketName string, client *minio.Client, bucketOptions minio.MakeBucketOptions) error {
	err := client.MakeBucket(context.Background(), bucketName, bucketOptions)
	if err != nil {
		return err
	}
	return nil
}

func DownloadObject(bucketName, objName, path string, client *minio.Client) error {
	err := client.FGetObject(context.Background(), bucketName, objName, path, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func PreSignedUrl(bucketName, objName string, client *minio.Client) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%v\"", objName))

	// Generates a presigned url which expires in an hour.
	presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, objName, time.Second*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// fmt.Print("Successfully generated presigned URL", presignedURL)
	return presignedURL, nil
}
