package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *CBAToolApp) callSummaryPage() {
	r := a.Data.CalcReport()
	summary := tview.NewTable().SetBorders(true)
	summary.SetTitle("Summary").SetBorder(true)
	summary.SetBordersColor(tcell.ColorYellow)
	summary.SetSelectable(false, false)
	// HEADER
	summary.SetCell(0, 0, tview.NewTableCell("").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		summary.SetCell(0, 1 + i, tview.NewTableCell(c.Name).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	}
	summary.SetCell(0, 1 + len(a.Data.Phases), tview.NewTableCell("Total").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	// end of HEADER
	// External costs
	summary.SetCell(1, 0, tview.NewTableCell("External").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		tot := r.External.Labor[c.Index] + r.External.Investment[c.Index] + r.External.Consulting[c.Index] + r.External.Others[c.Index]
		summary.SetCell(1, 1 + i, tview.NewTableCell(fmt.Sprintf("%.2f %s", tot, a.Data.Currency)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	}
	grandtot := r.External.TotLabor + r.External.TotInvestment + r.External.TotConsulting + r.External.TotOthers
	summary.SetCell(1, 1 + len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", grandtot, a.Data.Currency)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	// end of External costs
	// Internal costs
	summary.SetCell(2, 0, tview.NewTableCell("Internal").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		tot := r.Internal.Labor[c.Index] + r.Internal.Investment[c.Index] + r.Internal.Consulting[c.Index] + r.Internal.Others[c.Index]
		summary.SetCell(2, 1 + i, tview.NewTableCell(fmt.Sprintf("%.2f %s", tot, a.Data.Currency)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	}
	grandtot = r.Internal.TotLabor + r.Internal.TotInvestment + r.Internal.TotConsulting + r.Internal.TotOthers
	summary.SetCell(2, 1 + len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", grandtot, a.Data.Currency)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	// end of Internal costs
	a.app.SetRoot(summary, true)
}