package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *KirinukiApp) callQuit() {
	modal := tview.NewModal().SetText("Do you want to quit the application?").AddButtons([]string{"Quit", "Cancel"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			a.app.Stop()
		} else {
			a.callMenuPage()
		}
	})
	a.callPage(modal)
}

func (a *KirinukiApp) callMessage(message string, err error, callback func()) {
	modal := tview.NewModal()
	if err != nil {
		modal.SetTitle("Error message")
		modal.SetText(fmt.Sprintf("%s ( %v )", message, err))
	} else {
		modal.SetTitle("Message")
		modal.SetText(message)
	}
	modal.AddButtons([]string{"Continue"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		callback()
	})
	a.app.SetRoot(modal, true)
}

func (a *KirinukiApp) callLineInput(message string) {
	inputField := tview.NewInputField()
	inputField.SetLabel(message)
	inputField.SetFieldWidth(10)
	inputField.SetAcceptanceFunc(tview.InputFieldInteger)
	inputField.SetDoneFunc(func(key tcell.Key) {
		a.callMenuPage()
	})
}
