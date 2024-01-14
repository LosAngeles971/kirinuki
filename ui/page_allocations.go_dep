package ui

import (
	"fmt"

	"github.com/LosAngeles971/cba-tool/business/cba"
	"github.com/rivo/tview"
)

func (a *CBAToolApp) callAllocationsPage() {
	allocationsPage := tview.NewTable().SetBorders(true)
	allocationsPage.SetBorder(true).SetTitle("Allocations of costs to project's cycles")
	allocationsPage.SetSelectable(true, false)
	allocationsPage.SetSelectedFunc(func(row int, column int) {
		if row < 1 {
			a.callAddNewAllocation()
		} else {
			a.callUpdateAllocation(a.Data.Costs[row-1].Name)
		}
	})
	allocationsPage.Clear()
	allocationsPage.SetCell(0, 0, tview.NewTableCell("Cost").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	allocationsPage.SetCell(0, 1, tview.NewTableCell("Item occurrences").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	allocationsPage.SetCell(0, 2, tview.NewTableCell("Allocated to").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	allocationsPage.SetCell(0, 3, tview.NewTableCell("Applied discount").SetTextColor(t_colors[table_header_color]).SetAlign(tview.AlignCenter))
	if a.Data == nil {
		for i, alloc := range a.Data.Allocations {
			allocationsPage.SetCell(i+1, 0, tview.NewTableCell(alloc.Cost).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
			allocationsPage.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprint(alloc.Occurrence)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
			allocationsPage.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprint(alloc.Phase)).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
			allocationsPage.SetCell(i+1, 3, tview.NewTableCell(alloc.Discount).SetTextColor(t_colors[table_cell_color]).SetAlign(tview.AlignCenter))
		}
	}
	a.callPage(allocationsPage)
}

func (a *CBAToolApp) getAllocationForm(alloc *cba.Allocation, update bool) *tview.Form {
	form := tview.NewForm()
	form.AddDropDown("Associated cost", a.Data.ListCosts(), a.getCostIndex(alloc.Cost), nil)
	form.AddInputField("Cost occurrences", fmt.Sprint(alloc.Occurrence), 15, tview.InputFieldInteger, nil)
	form.AddDropDown("Associated phase", a.Data.ListPhases(), alloc.Phase, nil)
	form.AddButton("Add/Update", func() {
		_, c := form.GetFormItemByLabel("Associated cost").(*tview.DropDown).GetCurrentOption()
		_, n := form.GetFormItemByLabel("Associated phase").(*tview.DropDown).GetCurrentOption()
		p := a.Data.FindPhaseByName(n)
		if p != nil {
			alloc.Cost = c
			alloc.Phase = p.Index
			alloc.Occurrence = getFloat(form.GetFormItemByLabel("Cost occurrences").(*tview.InputField))
			a.callAllocationsPage()
		}
	})
	form.AddButton("Cancel", func() {
		a.callAllocationsPage()
	})
	if update {
		form.AddButton("Delete", func() {
			a.Data.DeleteAllocationByID(alloc.ID)
			a.callAllocationsPage()
		})
	}
	return form
}

func (a *CBAToolApp) callUpdateAllocation(id string) {
	alloc := a.Data.FindAllocationByID(id)
	if alloc == nil {
		return
	}
	form := a.getAllocationForm(alloc, true)
	form.SetBorder(true).SetTitle("Update cost allocation").SetTitleAlign(tview.AlignLeft)
	a.app.SetRoot(form, true)
}

func (a *CBAToolApp) callAddNewAllocation() {
	alloc := a.Data.NewAllocation()
	form := a.getAllocationForm(alloc, false)
	form.SetBorder(true).SetTitle("Add cost allocation").SetTitleAlign(tview.AlignLeft)
	a.app.SetRoot(form, true)
}
