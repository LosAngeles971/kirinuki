package ui

import (
	"fmt"

	"github.com/LosAngeles971/cba-tool/business/cba"
	"github.com/rivo/tview"
)

func (a *CBAToolApp) callPhasesPage() {
	phasesPage := tview.NewTable().SetBorders(true)
	phasesPage.SetBorder(true).SetTitle("Phases")
	phasesPage.SetSelectable(true, false)
	phasesPage.SetSelectedFunc(func(row int, column int) {
		if row < 1 {
			a.callAddNewPhase()
		} else {
			a.callUpdatePhase(a.Data.Phases[row-1].Index)
		}
	})
	phasesPage.Clear()
	phasesPage.SetCell(0, 0, tview.NewTableCell("Index").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	phasesPage.SetCell(0, 1, tview.NewTableCell("Phase").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	phasesPage.SetCell(0, 2, tview.NewTableCell("Days").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	if a.Data != nil {
		for i, cycle := range a.Data.Phases {
			phasesPage.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprint(cycle.Index)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
			phasesPage.SetCell(i+1, 1, tview.NewTableCell(cycle.Name).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
			phasesPage.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprint(cycle.Days)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
		}
	}
	a.app.SetRoot(phasesPage, true)
}

func (a *CBAToolApp) getPhaseForm(phase *cba.Phase, update bool) *tview.Form {
	form := tview.NewForm()
	form.AddInputField("Name", phase.Name, 40, nil, nil)
	form.AddInputField("Index", fmt.Sprint(phase.Index), 3, tview.InputFieldInteger, nil)
	form.AddInputField("Days", fmt.Sprint(phase.Days), 10, tview.InputFieldInteger, nil)
	form.AddButton("Add/Update", func() {
		phase.Name = form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
		phase.Index = getInteger(form.GetFormItemByLabel("Index").(*tview.InputField))
		phase.Days = getInteger(form.GetFormItemByLabel("Days").(*tview.InputField))
		if !update {
			a.Data.Phases = append(a.Data.Phases, phase)
		}
		a.callPhasesPage()
	})
	form.AddButton("Cancel", func() {
		a.callPhasesPage()
	})
	if update {
		form.AddButton("Delete", func() {
			a.Data.DeletePhaseByIndex(phase.Index)
			a.callPhasesPage()
		})
	}
	return form
}

func (a *CBAToolApp) callUpdatePhase(i int) {
	phase := a.Data.FindPhaseByIndex(i)
	if phase == nil {
		return
	}
	form := a.getPhaseForm(phase, true)
	form.SetBorder(true).SetTitle("Update phase").SetTitleAlign(tview.AlignLeft)
	a.app.SetRoot(form, true)
}

func (a *CBAToolApp) callAddNewPhase() {
	phase := a.Data.NewPhase()
	form := a.getPhaseForm(phase, false)
	form.SetBorder(true).SetTitle("Enter new project's phase").SetTitleAlign(tview.AlignLeft)
	a.app.SetRoot(form, true)
}
