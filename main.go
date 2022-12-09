package main

import (
	"log"

	"github.com/one2nc/minio-tui/minio"
	"github.com/one2nc/minio-tui/tui"
	"github.com/rivo/tview"
)

var minioCfg = &minio.Config{
	Endpoint:        "play.min.io",
	AccessKey:       "Q3AM3UQ867SPQQA43P2F",
	SecretAccessKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
	UseSSL:          true,
}

func main() {
	
	minioClinet, err := minio.GetMinioClient(minioCfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	buckets, err := minio.GetBuckets(minioClinet)
	if err != nil {
		log.Fatalf("%v", err)
	}

	app := tview.NewApplication()
	pages := tview.NewPages()
	tuiConfig := &tui.Config{
		App:         app,
		Pages:       pages,
		MinioClient: minioClinet,
	}

	flex := tui.DisplayBuckets(buckets, tuiConfig)
	tuiConfig.Pages.AddAndSwitchToPage("page1", flex, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
