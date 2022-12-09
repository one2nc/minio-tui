package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/minio/minio-go/v7"
	m "github.com/one2nc/minio-tui/minio"
	"github.com/rivo/tview"
)

type Config struct {
	App         *tview.Application
	Pages       *tview.Pages
	MinioClient *minio.Client
}

func DisplayBuckets(buckets []minio.BucketInfo, config *Config) (page *tview.Flex) {
	flex := tview.NewFlex()

	text := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText(fmt.Sprintf("Stats:\nBuckets: %v\nControls\nmouse click\nenter\n<?> help\n<r> Refresh Buckets\n<c> create new bucket", len(buckets)))
	text.SetBorder(true)

	table := tview.NewTable()
	table.SetBorder(true)

	//layout
	flex.AddItem(text, 0, 1, true).SetDirection(tview.FlexColumn)
	flex.AddItem(table, 0, 4, true).SetDirection(tview.FlexRow)

	//table data
	table.SetCell(0, 0, tview.NewTableCell("Bucket-Name").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("CreationDate").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	//table.SetCell(0, 2, tview.NewTableCell("Size").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))

	for i, b := range buckets {
		table.SetCell((i + 1), 0, tview.NewTableCell(b.Name).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(b.CreationDate.String()).SetAlign(tview.AlignCenter))
		//table.SetCell((i + 1), 2, tview.NewTableCell(fmt.Sprintf("%f", b.Size)).SetAlign(tview.AlignCenter))
	}

	table.Select(1, 1).SetFixed(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		config.Pages.RemovePage("page2")
		//pages.AddPage("page1", DisplayBuckets(buckets, pages, app), true, false)
		files, err := m.GetFiles(buckets[row].Name, config.MinioClient)
		if err != nil {
			fmt.Println("error: ", err)
		}
		config.Pages.AddPage("page2", DisplayFiles(files, config), true, false)

		config.Pages.SwitchToPage("page2")
		table.SetSelectable(true, false)
	})

	config.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 114 {
			//extract to methhod **Refresh Bucket**
			config.Pages.RemovePage("page1")
			buckets, err := m.GetBuckets(config.MinioClient)
			if err != nil {
				fmt.Println("error: ", err)
			}
			page := DisplayBuckets(buckets, config)
			config.Pages.AddAndSwitchToPage("page1", page, true)
		} else if event.Rune() == 99 {
			// **Create Bucket**
			// var bucketName string
			// // var form = tview.NewForm()
			// // form.AddInputField("Bucket Name: ", "", 20, nil, func(bucName string) {
			// // 	bucketName = bucName
			// // })

			// // SetText("Enter your name:").
			// // AddButtons([]string{"Submit"})

			// // modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// // 	if buttonLabel == "Submit" {
			// // 		// Get input from the modal and process it
			// // 		bucketName = modal.GetFormItemText("bucketName")
			// // 	}
			// // })
			// // modal.AddFormItem("bucketName", "Your name:", "", 0, nil)
			// modal := tview.NewModal().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// 	switch event.Key() {
			// 	case tcell.KeyEscape:
			// 		// Handle the escape key
			// 	case tcell.KeyEnter:
			// 		// Handle the enter key
			// 	default:
			// 		// Handle other keys
			// 	}

			// 	return event
			// })
			// config.Pages.AddAndSwitchToPage("form", modal, true)

			// m.MakeBucket(bucketName, config.MinioClient, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})

			// 	SetText(fmt.Sprintf("File Info\nName: %v\nOwner: %v\nLastModified: %v\nExpiration: %v\nHasDeleteMarker: %v", file.Key, file.Owner.DisplayName, file.LastModified, file.Expiration, file.IsDeleteMarker)).
			// 	AddButtons([]string{"Ok"}).
			// 	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// 		config.Pages.SwitchToPage("page2")
			// 	})
			// modal.SetBackgroundColor(tcell.ColorDarkBlue)
			// config.Pages.AddAndSwitchToPage("modal", modal, true)
			// config.Pages.RemovePage("form")

			// buckets, err := m.GetBuckets(config.MinioClient)
			// if err != nil {
			// 	fmt.Println("error: ", err)
			// }
			// page := DisplayBuckets(buckets, config)
			// config.Pages.AddAndSwitchToPage("page1", page, true)

		}
		return event
	})

	return flex
}

type File struct {
	Name    string
	LastMod time.Time
	Size    float32
}

func DisplayFiles(files []minio.ObjectInfo, config *Config) *tview.Flex {

	text := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("Page 2 Controls\n <b> back(<-)")
	text.SetBorder(true)

	table := tview.NewTable()
	table.SetBorder(true)

	flex := tview.NewFlex()
	//layout
	flex.AddItem(text, 0, 1, true).SetDirection(tview.FlexColumn)
	flex.AddItem(table, 0, 4, true).SetDirection(tview.FlexRow)

	//table data
	table.SetCell(0, 0, tview.NewTableCell("File-Name").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("IsLatest").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("ReplicationReady").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("Last Modified").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	table.SetCell(0, 4, tview.NewTableCell("Size").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))

	for i, f := range files {
		table.SetCell((i + 1), 0, tview.NewTableCell(f.Key).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(fmt.Sprintf("%v", f.IsLatest)).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 2, tview.NewTableCell(fmt.Sprintf("%v", f.ReplicationReady)).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 3, tview.NewTableCell(f.LastModified.GoString()).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 4, tview.NewTableCell(fmt.Sprintf("%v", f.Size)).SetAlign(tview.AlignCenter))
	}

	table.Select(1, 1).SetFixed(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.SetSelectable(true, false)
		file := files[row-1]
		modal := tview.NewModal().
			SetText(fmt.Sprintf("File Info\nName: %v\nOwner: %v\nLastModified: %v\nExpiration: %v\nHasDeleteMarker: %v", file.Key, file.Owner.DisplayName, file.LastModified, file.Expiration, file.IsDeleteMarker)).
			AddButtons([]string{"Ok"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				config.Pages.SwitchToPage("page2")
			})
		modal.SetBackgroundColor(tcell.ColorDarkBlue)
		config.Pages.AddAndSwitchToPage("modal", modal, true)
	})

	//capture events
	config.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 98 {
			config.Pages.SwitchToPage("page1")
		}
		return event
	})

	return flex
}
