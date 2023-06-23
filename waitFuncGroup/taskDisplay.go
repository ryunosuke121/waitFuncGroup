package waitFuncGroup

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application

func createTable() (table *tview.Table) {
	newtable := tview.NewTable().SetBorders(true)

	newtable.SetCell(0, 0, tview.NewTableCell("task").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignCenter))
	newtable.SetCell(0, 1, tview.NewTableCell("status").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignCenter))
	return newtable
}

func display(table *tview.Table) {
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func setCompleteRow(table *tview.Table, taskid int, name string) {
	app.QueueUpdateDraw(func() {
		table.SetCell(taskid, 0, tview.NewTableCell(name).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(taskid, 1, tview.NewTableCell("completed").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	})
}

func setWorkingRow(table *tview.Table, taskid int, name string) {
	app.QueueUpdateDraw(func() {
		table.SetCell(taskid, 0, tview.NewTableCell(name).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(taskid, 1, tview.NewTableCell("working...").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	})
}

func setPanicRow(table *tview.Table, taskid int, name string) {
	app.QueueUpdateDraw(func() {
		table.SetCell(taskid, 0, tview.NewTableCell(name).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(taskid, 1, tview.NewTableCell("panic!").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	})
}
