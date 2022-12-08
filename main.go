package main

import (
	"time"

	"github.com/one2nc/minio-tui/tui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	data := []tui.Bucket{{BucName: "test-bucket1", LastMod: time.Now(), Size: 30}, {BucName: "test-bucket2", LastMod: time.Now(), Size: 10330}, {BucName: "test-bucket3", LastMod: time.Now(), Size: 30}, {BucName: "test-bucket4", LastMod: time.Now(), Size: 100230}}
	flex := tui.DisplayBuckets(data)
	if err := app.SetRoot(&flex, true).EnableMouse(true).SetFocus(&flex).Run(); err != nil {
		panic(err)
	}
}
