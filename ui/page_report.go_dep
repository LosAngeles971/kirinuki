package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *CBAToolApp) callReportPage() {
	r := a.Data.CalcReport()
	reportPage := tview.NewTable().SetBorders(true)
	reportPage.SetBorder(true).SetTitle("Cost-Benefit Analysis report")
	reportPage.SetSelectable(false, false)
	// HEADER
	reportPage.SetCell(0, 0, tview.NewTableCell("").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(0, 1+i, tview.NewTableCell(c.Name).SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(0, 1+len(a.Data.Phases), tview.NewTableCell("Total").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
	// end of HEADER
	// External LABOR
	reportPage.SetCell(1, 0, tview.NewTableCell("External - labor").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(1, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.Labor[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(1, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.TotLabor, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of External LABOR
	// External INVESTMENT
	reportPage.SetCell(2, 0, tview.NewTableCell("External - investment").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(2, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.Investment[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(2, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.TotInvestment, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of External INVESTMENT
	// External CONSULTING
	reportPage.SetCell(3, 0, tview.NewTableCell("External - consulting").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(3, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.Consulting[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(3, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.TotConsulting, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of External CONSULTING
	// External OTHERS
	reportPage.SetCell(4, 0, tview.NewTableCell("External - others").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(4, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.Others[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(4, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.External.TotOthers, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of External OTHERS
	// Internal LABOR
	reportPage.SetCell(5, 0, tview.NewTableCell("Internal - labor").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(5, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.Labor[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(5, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.TotLabor, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of Internal LABOR
	// Internal INVESTMENT
	reportPage.SetCell(6, 0, tview.NewTableCell("Internal - investment").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(6, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.Investment[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(6, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.TotInvestment, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of Internal INVESTMENT
	// Internal CONSULTING
	reportPage.SetCell(7, 0, tview.NewTableCell("Internal - consulting").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(7, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.Consulting[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(7, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.TotConsulting, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of Internal CONSULTING
	// Internal OTHERS
	reportPage.SetCell(8, 0, tview.NewTableCell("Internal - others").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	for i, c := range a.Data.Phases {
		reportPage.SetCell(8, 1+i, tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.Others[c.Index], a.Data.Currency)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
	}
	reportPage.SetCell(8, 1+len(a.Data.Phases), tview.NewTableCell(fmt.Sprintf("%.2f %s", r.Internal.TotOthers, a.Data.Currency)).SetTextColor(t_colors[table_highcell_color]).SetAlign(tview.AlignCenter))
	// end of Internal OTHERS
	a.app.SetRoot(reportPage, true)
}
