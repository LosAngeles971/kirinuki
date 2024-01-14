package ui

import (
	"fmt"
	"strconv"

	"github.com/LosAngeles971/cba-tool/business/cba"
	"github.com/rivo/tview"
)

func (a *CBAToolApp) callUpdateSettings() {
	var form *tview.Form
	ic := 0
	if a.Data.Currency != cba.CURRENCY_EURO {
		ic = 1
	}
	form = tview.NewForm().AddDropDown("Currency", []string{cba.CURRENCY_EURO, cba.CURRENCY_DOLLAR,}, ic, nil).
		AddInputField("VAT (%)", fmt.Sprint(a.Data.ValueAddedTax), 5, tview.InputFieldFloat, nil).
		AddButton("Confirm", func() {
			_, c := form.GetFormItemByLabel("Currency").(*tview.DropDown).GetCurrentOption()
			a.Data.Currency = c
			vat, _ := strconv.ParseFloat(form.GetFormItemByLabel("Name").(*tview.InputField).GetText(), 64)
			a.Data.ValueAddedTax = vat
			a.callMenuPage()
		}).AddButton("Cancel", a.callMenuPage)
	form.SetBorder(true).SetTitle("Project's settings").SetTitleAlign(tview.AlignLeft)
	a.app.SetRoot(form, true)
}
