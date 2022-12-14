package main

import (
	"log"
	// x "github.com/minio/minio-go/v7"
	"github.com/one2nc/minio-tui/minio"
	"github.com/one2nc/minio-tui/tui"
	"github.com/rivo/tview"
)

var minioCfg = &minio.Config{
	// Endpoint:        "play.min.io",
	// AccessKey:       "Q3AM3UQ867SPQQA43P2F",
	// SecretAccessKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",

	Endpoint:        "127.0.0.1:9000",
	AccessKey:       "minioadmin",
	SecretAccessKey: "minioadmin",
	UseSSL:          false,
}

func main() {
	
	minioClient, err := minio.GetMinioClient(minioCfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

    // bucketName:="one2nbucket"
	// minio.MakeBucket(bucketName,minioClient,x.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
	buckets, err := minio.GetBuckets(minioClient)
	if err != nil {
		log.Fatalf("%v", err)
	}

	app := tview.NewApplication()
	pages := tview.NewPages()
	tuiConfig := &tui.Config{
		App:         app,
		Pages:       pages,
		MinioClient: minioClient,
	}

	flex := tui.DisplayBuckets(buckets, tuiConfig)
	tuiConfig.Pages.AddAndSwitchToPage("page1", flex, true)

	if err := app.SetRoot(pages, true).EnableMouse(false).SetFocus(pages).Run(); err != nil {
		panic(err)
	}


}
