package ui

/*
https://github.com/rivo/tview/wiki/Table
https://gist.github.com/rivo/2893c6740a6c651f685b9766d1898084
https://github.com/rivo/tview/wiki/Postgres
https://github.com/rivo/tview/wiki
https://github.com/destinmoulton/pixi/blob/master/gui/gui.go
https://github.com/rivo/tview
*/
import (
	"github.com/LosAngeles971/kirinuki/business"
	"github.com/LosAngeles971/kirinuki/business/storage"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KirinukiApp struct {
	gw  *business.Gateway
	ms  *storage.MultiStorage
	app *tview.Application
}

func (a *KirinukiApp) eventHandler(eventKey *tcell.EventKey) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEscape {
		a.callMenuPage()
	}
	return eventKey
}

func Build() *KirinukiApp {
	kk := &KirinukiApp{
		gw: nil,
		ms: nil,
	}
	kk.app = tview.NewApplication()

	kk.app.EnableMouse(true)
	kk.app.SetInputCapture(kk.eventHandler)

	kk.callMenuPage()

	return kk
}

func (a *KirinukiApp) Run() error {
	return a.app.Run()
}
