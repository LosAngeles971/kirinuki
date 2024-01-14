package ui

import (
	"github.com/LosAngeles971/kirinuki/business"
	"github.com/rivo/tview"
)

func (a *KirinukiApp) callOpenTOC() {
	var form *tview.Form
	if a.gw != nil {
		a.callMessage("You must logout from current Table of Content before", nil, a.callMenuPage)
	}
	if a.ms == nil {
		a.callMessage("You must load the storage map before", nil, a.callMenuPage)
	}
	form = tview.NewForm().
		AddInputField("Email", "", 40, nil, nil).
		AddPasswordField("Passphrase", "", 40, '*', nil).
		AddButton("Open", func() {
			var err error
			a.gw, err = business.New(
				form.GetFormItemByLabel("Email").(*tview.InputField).GetText(), 
				form.GetFormItemByLabel("Passphrase").(*tview.InputField).GetText(), 
				business.WithStorage(a.ms))
			if err != nil {
				a.callMessage("failed to open the table of content", err, a.callMenuPage)
			} else {
				a.callMenuPage()
			}
		}).AddButton("Cancel", a.callMenuPage)
	form.SetBorder(true).SetTitle("Project's settings").SetTitleAlign(tview.AlignLeft)
	a.app.SetRoot(form, true)
}
