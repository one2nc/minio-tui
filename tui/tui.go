package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/minio/minio-go/v7"
	m "github.com/one2nc/minio-tui/minio"
	"github.com/rivo/tview" //replaced with "github.com/derailed/tview v0.7.2"
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
		config.Pages.AddPage("page2", DisplayFiles(buckets[row].Name, files, config), true, false)

		config.Pages.SwitchToPage("page2")
		table.SetSelectable(true, false)
	})

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			cbmf, _ := makeCreateBucketModalForm(config, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: false})
			config.Pages.RemovePage("create-bucket")
			config.Pages.AddPage("create-bucket", cbmf, true, true)
		}
		return event
	})

	return flex
}

func makeCreateBucketModalForm(config *Config, bucketOpt minio.MakeBucketOptions) (*tview.ModalForm, error) {
	// **Create Bucket**
	bktName := ""
	bktForm := tview.NewForm()
	errorTxt := tview.NewTextView()
	errorTxt.SetTextColor(tcell.ColorIndianRed)

	bktForm.AddInputField("Bucket Name: ", "", 20, nil, func(bucketName string) {
		bktName = bucketName
	})
	bktForm.SetFieldBackgroundColor(tcell.ColorBlack.TrueColor())
	var err error
	createBktModal := tview.NewModalForm("<Create Bucket>", bktForm)
	createBktModal.SetBorder(true).SetBorderColor(tcell.ColorBlue)
	createBktModal.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	createBktModal.AddButtons([]string{"create", "cancel"})
	createBktModal.SetButtonBackgroundColor(tcell.ColorDarkRed)
	createBktModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "create" {
			err = m.MakeBucket(bktName, config.MinioClient, bucketOpt)
			if err == nil {
				config.Pages.RemovePage("page1")
				buckets, _ := m.GetBuckets(config.MinioClient)
			    config.Pages.AddPage("page1", DisplayBuckets(buckets, config), true, false)
				makeAckModel(config, "page1", "Success", "Bucket Created!")
			}
			if err != nil {
				makeAckModel(config, "create-bucket", "Error", err.Error())
			}
		}
		if buttonLabel == "cancel" {
			config.Pages.SwitchToPage("page1")
		}
	})
	return createBktModal, err
}

func makeAckModel(config *Config, currentPage, title, msg string) {
	modal := tview.NewModal()
	modal.SetTitle(title)
	modal.SetText(msg)
	modal.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	modal.AddButtons([]string{"OK"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		config.Pages.RemovePage("modal")
		config.Pages.SwitchToPage(currentPage)
	})
	config.Pages.AddPage("modal", modal, true, true)
}

func DisplayFiles(bucketName string, files []minio.ObjectInfo, config *Config) *tview.Flex {
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

	r := 0
	table.Select(1, 1).SetFixed(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter || key == tcell.KeyUp || key == tcell.KeyDown {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.SetSelectable(true, false)
		// file := files[row-1]
		r = row - 1

		// modal := tview.NewModal().
		// 	SetText(fmt.Sprintf("File Info\nName: %v\nOwner: %v\nLastModified: %v\nExpiration: %v\nHasDeleteMarker: %v", file.Key, file.Owner.DisplayName, file.LastModified, file.Expiration, file.IsDeleteMarker)).
		// 	AddButtons([]string{"Ok"}).
		// 	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		// 		config.Pages.SwitchToPage("page2")
		// 	})
		// modal.SetBackgroundColor(tcell.ColorDarkBlue)
		// config.Pages.AddAndSwitchToPage("modal", modal, true)
	})

	//capture events
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 98 {
			config.Pages.SwitchToPage("page1")
		}
		if event.Key() == tcell.KeyCtrlD {
			err := m.DownloadObject(bucketName, files[r].Key, "../resources/downloads", config.MinioClient)
			if err != nil {
				makeAckModel(config, "page2", "<Error>", err.Error())
			}else{
				makeAckModel(config, "page2", "<Greet>", fmt.Sprintf("%v downloaded.", files[r].Key))
			}
		}
		return event
	})

	return flex
}
