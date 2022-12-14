package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/minio/minio-go/v7"
	m "github.com/one2nc/minio-tui/minio"
	"github.com/rivo/tview" //replaced with "github.com/derailed/tview v0.7.2"
	"github.com/sahilm/fuzzy"
)

type Config struct {
	App         *tview.Application
	Pages       *tview.Pages
	MinioClient *minio.Client
}

type Prompt struct {
	*tview.TextView
	noIcons bool
	icon    rune
	spacer  int
}

func DisplayBuckets(buckets []minio.BucketInfo, config *Config) *tview.Flex {
	page := tview.NewFlex()
	//header
	header := tview.NewFlex()
	stats := tview.NewTextView().
		//SetTextColor(tcell.Color(tcell.ColorDarkGreen)).
		SetText(fmt.Sprintf("Buckets: %v\n", len(buckets)))
	stats.SetBorderPadding(1, 1, 1, 1)

	controls := tview.NewTextView().SetText("<?> help\n<r> Refresh Buckets\n<c> create new bucket\n</> Search")
	controls.SetBorderPadding(1, 1, 1, 1)

	header.AddItem(stats, 0, 5, false).SetDirection(tview.FlexColumn)
	header.AddItem(controls, 0, 5, false).SetDirection(tview.FlexColumn)
	//controls.SetBorder(true)

	//content
	table := tview.NewTable()
	table.SetBorderPadding(1, 1, 1, 1)
	table.SetBorder(true).SetBorderColor(tcell.NewRGBColor(205, 133, 63))
	table.SetTitle(" BUCKETS ").SetTitleColor(tcell.NewRGBColor(139, 69, 19))

	//footer: for acknowlege things
	ack := tview.NewTextView()

	//ack.SetBorder(true)
	ack.SetTextAlign(tview.AlignCenter)
	ack.SetText("")

	//layout
	page.AddItem(header, 0, 3, false).SetDirection(tview.FlexColumn)
	page.AddItem(table, 0, 8, true).SetDirection(tview.FlexRow)
	page.AddItem(ack, 0, 2, false).SetDirection(tview.FlexRow)
	//table data
	h1 := tview.NewTableCell("NAME").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter)
	table.SetCell(0, 0, h1)
	table.SetCell(0, 1, tview.NewTableCell("CREATION DATE").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	//table.SetCell(0, 2, tview.NewTableCell("Size").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))

	table.SetFixed(1, 1)
	for i, b := range buckets {
		table.SetCell((i + 1), 0, tview.NewTableCell(b.Name).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(b.CreationDate.String()).SetAlign(tview.AlignCenter))
		//table.SetCell((i + 1), 2, tview.NewTableCell(fmt.Sprintf("%f", b.Size)).SetAlign(tview.AlignCenter))
	}

	table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	})

	table.SetSelectedFunc(func(row int, column int) {
		config.Pages.RemovePage("page2")
		//pages.AddPage("page1", DisplayBuckets(buckets, pages, app), true, false)
		row = row - 1
		if row < 0 {
			row = 0
		}
		files, err := m.GetFiles(buckets[row].Name, config.MinioClient)
		if err != nil {
			ack.SetText(err.Error())
		}
		config.Pages.AddPage("page2", DisplayFiles(buckets[row].Name, files, config), true, false)

		config.Pages.SwitchToPage("page2")
		table.SetSelectable(true, false)
	})

	page.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 114 {
			//extract to methhod **Refresh Bucket**
			config.Pages.RemovePage("page1")
			buckets, err := m.GetBuckets(config.MinioClient)
			if err != nil {
				ack.SetText(err.Error())
			}
			page := DisplayBuckets(buckets, config)
			config.Pages.AddAndSwitchToPage("page1", page, true)
		} else if event.Rune() == 99 {
			cbmf, _ := makeCreateBucketModalForm(config, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: false})
			config.Pages.RemovePage("create-bucket")
			config.Pages.AddPage("create-bucket", cbmf, true, true)
		} else if event.Rune() == 47 {
			form, _ := FilterBucketForm(config, buckets)
			config.Pages.RemovePage("search-filter")
			config.Pages.AddPage("search-filter", form, true, true)
		}
		return event
	})

	//frame := tview.NewFrame(page)
	return page
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
		SetText("Page 2 Controls\n <b> back(<-)\n<ctrl + d> Download File\n</> Search\n<ctrl + p> Presigned url")
	text.SetBorder(true)

	table := tview.NewTable()
	table.SetBorder(true)

	flex := tview.NewFlex()
	//layout
	flex.AddItem(text, 0, 1, false).SetDirection(tview.FlexColumn)
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
	table.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.SetSelectable(true, false)
		r = row - 1
	})

	//capture events
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//get back
		if event.Rune() == 98 {
			config.Pages.SwitchToPage("page1")
		} else if event.Rune() == 47 {
			form, _ := FilterFileForm(bucketName, config, files)
			config.Pages.RemovePage("search-file")
			config.Pages.AddPage("search-file", form, true, true)
		}

		//download selected object
		if event.Key() == tcell.KeyCtrlD {
			objName := files[r].Key
			path := fmt.Sprintf("./resources/downloads/%v", objName)
			err := m.DownloadObject(bucketName, objName, path, config.MinioClient)
			if err != nil {
				makeAckModel(config, "page2", "<Error>", err.Error())
			} else {
				makeAckModel(config, "page2", "<Greet>", fmt.Sprintf("%v downloaded.", files[r].Key))
			}
		}
		if event.Key() == tcell.KeyCtrlP {
			p, err := m.PreSignedUrl(bucketName, files[r].Key, config.MinioClient)
			if err != nil {
				makeAckModel(config, "page2", "<Error>", err.Error())
			} else {
				makeAckModel(config, "page2", "<Presigned URL>", fmt.Sprintf("%v", p))
			}
		}
		return event
	})

	return flex
}

func SearchBucketFilter(searchString string, buckets []minio.BucketInfo) []minio.BucketInfo {
	var bucketsTemp []string
	for _, bucket := range buckets {
		bucketsTemp = append(bucketsTemp, bucket.Name)
	}
	match := fuzzy.Find(searchString, bucketsTemp)
	var bucketsTempWT []minio.BucketInfo
	for i := 0; i < len(match); i++ {
		for j := 0; j < len(buckets); j++ {
			if strings.EqualFold(buckets[j].Name, match[i].Str) {
				bucketsTempWT = append(bucketsTempWT, buckets[j])
			}
		}
	}
	return bucketsTempWT
}

func SearchFileFilter(searchString string, files []minio.ObjectInfo) []minio.ObjectInfo {
	var bucketsTemp []string
	for _, file := range files {
		bucketsTemp = append(bucketsTemp, file.Key)
	}
	match := fuzzy.Find(searchString, bucketsTemp)
	var filesTempWT []minio.ObjectInfo
	for i := 0; i < len(match); i++ {
		for j := 0; j < len(files); j++ {
			if strings.EqualFold(files[j].Key, match[i].Str) {
				filesTempWT = append(filesTempWT, files[j])
			}
		}
	}
	return filesTempWT
}

func FilterBucketForm(config *Config, buckets []minio.BucketInfo) (*tview.ModalForm, error) {
	searchBucket := ""
	bktForm := tview.NewForm()
	errorTxt := tview.NewTextView()
	errorTxt.SetTextColor(tcell.ColorIndianRed)

	bktForm.AddInputField("Search: ", "", 20, nil, func(bucketName string) {
		searchBucket = bucketName
	})
	bktForm.SetFieldBackgroundColor(tcell.ColorBlack.TrueColor())
	var err error
	filterBucketModal := tview.NewModalForm("<Filter Form>", bktForm)
	filterBucketModal.SetBorder(true).SetBorderColor(tcell.ColorBlue)
	filterBucketModal.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	filterBucketModal.AddButtons([]string{"search", "cancel"})
	filterBucketModal.SetButtonBackgroundColor(tcell.ColorDarkRed)
	filterBucketModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "search" {
			bucketsTempWT := SearchBucketFilter(searchBucket, buckets)
			config.Pages.RemovePage("filterBucketsPage")
			page := DisplayFilterBuckets(bucketsTempWT, config)
			config.Pages.AddAndSwitchToPage("filterBucketsPage", page, true)
		}
		if buttonLabel == "cancel" {
			config.Pages.SwitchToPage("page1")
		}
	})
	return filterBucketModal, err
}

func FilterFileForm(bucketName string, config *Config, files []minio.ObjectInfo) (*tview.ModalForm, error) {
	searchFile := ""
	bktForm := tview.NewForm()
	errorTxt := tview.NewTextView()
	errorTxt.SetTextColor(tcell.ColorIndianRed)

	bktForm.AddInputField("Search: ", "", 20, nil, func(fileName string) {
		searchFile = fileName
	})
	bktForm.SetFieldBackgroundColor(tcell.ColorBlack.TrueColor())
	var err error
	filterFormModal := tview.NewModalForm("<Filter Form>", bktForm)
	filterFormModal.SetBorder(true).SetBorderColor(tcell.ColorBlue)
	filterFormModal.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	filterFormModal.AddButtons([]string{"search", "cancel"})
	filterFormModal.SetButtonBackgroundColor(tcell.ColorDarkRed)
	filterFormModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "search" {
			filesTempWT := SearchFileFilter(searchFile, files)
			config.Pages.RemovePage("filterFilesPage")
			page := DisplayFilterFiles(bucketName, filesTempWT, config)
			config.Pages.AddAndSwitchToPage("filterFilesPage", page, true)
		}
		if buttonLabel == "cancel" {
			config.Pages.SwitchToPage("page1")
		}
	})
	return filterFormModal, err
}

func DisplayFilterBuckets(buckets []minio.BucketInfo, config *Config) (page *tview.Flex) {
	flex := tview.NewFlex()

	text := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText(fmt.Sprintf("Stats:\nBuckets: %v\nControls\nMouse click\nEnter\n<?> Help\n<r> Refresh Buckets\n<c> Create new bucket\n <b> Back", len(buckets)))
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
		} else if event.Rune() == 47 {
			buckets, err := m.GetBuckets(config.MinioClient)
			if err != nil {
				fmt.Println("error: ", err)
			}
			form, _ := FilterBucketForm(config, buckets)
			config.Pages.RemovePage("search-bucket")
			config.Pages.AddPage("search-bucket", form, true, true)
		} else if event.Rune() == 98 {
			config.Pages.RemovePage("page1")
			buckets, err := m.GetBuckets(config.MinioClient)
			if err != nil {
				fmt.Println("error: ", err)
			}
			page := DisplayBuckets(buckets, config)
			config.Pages.AddAndSwitchToPage("page1", page, true)
		}
		return event
	})
	return flex
}

func DisplayFilterFiles(bucketName string, files []minio.ObjectInfo, config *Config) *tview.Flex {
	text := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("Page 2 Controls\n <b> back(<-) \n </>Search")
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
		r = row - 1
	})

	//capture events
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 98 {
			config.Pages.SwitchToPage("page1")
		} else if event.Rune() == 47 {
			files, err := m.GetFiles(bucketName, config.MinioClient)
			if err != nil {
				fmt.Println("error: ", err)
			}
			form, _ := FilterFileForm(bucketName, config, files)
			config.Pages.RemovePage("search-file")
			config.Pages.AddPage("search-file", form, true, true)
		}
		if event.Key() == tcell.KeyCtrlD {
			err := m.DownloadObject(bucketName, files[r].Key, "../resources/downloads", config.MinioClient)
			if err != nil {
				makeAckModel(config, "page2", "<Error>", err.Error())
			} else {
				makeAckModel(config, "page2", "<Greet>", fmt.Sprintf("%v downloaded.", files[r].Key))
			}
		}
		return event
	})
	return flex
}
