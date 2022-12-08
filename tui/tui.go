package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Bucket struct {
	BucName string
	LastMod time.Time
	Size    float32
}

func Hello() string {
	return "hello"
}

func DisplayBuckets(buckets []Bucket) tview.Flex {
	flex := tview.NewFlex()

	text := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("Controls\n<?> help\nenter\nselect\nmouse click")
	text.SetBorder(true)

	table := tview.NewTable()
	table.SetBorder(true)

	//layout
	flex.AddItem(text, 0, 1, true).SetDirection(tview.FlexColumn)
	flex.AddItem(table, 0, 4, true).SetDirection(tview.FlexRow)

	//table data
	table.SetCell(0, 0, tview.NewTableCell("Bucket-Name").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Last Modified").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Size").SetTextColor(tcell.ColorDarkOliveGreen).SetAlign(tview.AlignCenter))

	for i, b := range buckets {
		table.SetCell((i + 1), 0, tview.NewTableCell(b.BucName).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 1, tview.NewTableCell(b.LastMod.String()).SetAlign(tview.AlignCenter))
		table.SetCell((i + 1), 2, tview.NewTableCell(fmt.Sprintf("%f", b.Size)).SetAlign(tview.AlignCenter))
	}

	table.Select(1, 1).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		text.SetText(Hello() + " " + buckets[row-1].BucName + buckets[row-1].LastMod.GoString())
		table.GetCell(row, column)
		table.SetSelectable(true, false)
	})

	return *flex
}
