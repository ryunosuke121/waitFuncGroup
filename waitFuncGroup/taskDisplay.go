package waitFuncGroup

import (
	"strconv"

	"github.com/rivo/tview"
)

var app *tview.Application

func createTable() (table *tview.Table) {
	newtable := tview.NewTable()

	newtable.SetCell(0, 0, tview.NewTableCell("task").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	newtable.SetCell(0, 1, tview.NewTableCell("status").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	return newtable
}

func display(table *tview.Table) {
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func setCompleteRow(table *tview.Table, taskid int) {
	app.QueueUpdateDraw(func() {
		table.SetCell(taskid, 0, tview.NewTableCell(strconv.Itoa(taskid)).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(taskid, 1, tview.NewTableCell("completed").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	})
}

func setWorkingRow(table *tview.Table, taskid int) {
	app.QueueUpdateDraw(func() {
		table.SetCell(taskid, 0, tview.NewTableCell(strconv.Itoa(taskid)).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(taskid, 1, tview.NewTableCell("working...").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	})
}

func setPanicRow(table *tview.Table, taskid int) {
	app.QueueUpdateDraw(func() {
		table.SetCell(taskid, 0, tview.NewTableCell(strconv.Itoa(taskid)).SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
		table.SetCell(taskid, 1, tview.NewTableCell("panic!").SetTextColor(tview.Styles.PrimaryTextColor).SetAlign(tview.AlignCenter))
	})
}
