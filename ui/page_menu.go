package ui

import (
	_ "embed"

	"github.com/rivo/tview"
)

//go:embed about.txt
var about_text string

//go:embed help.txt
var help_text string

func (a *KirinukiApp) callMenuPage() {
	flex := tview.NewFlex()
	help := tview.NewTextView()
	help.SetText(about_text)
	help.SetTextAlign(tview.AlignLeft)
	help.SetTitle("Help")
	help.SetBorder(true)
	help.SetWordWrap(true)
	help.SetDynamicColors(true)
	mainMenu := tview.NewList().ShowSecondaryText(false)
	mainMenu.SetBorder(true).SetTitle("Main menù")
	mainMenu.AddItem("About", "", ' ', func() {
		help.SetText(about_text)
	})
	mainMenu.AddItem("Help", "", ' ', func() {
		help.SetText(help_text)
	})
	mainMenu.AddItem("Open Table of Content", "", 'T', func() {
		a.callOpenTOC()
	})
	mainMenu.AddItem("Quit", "", 'Q', func() {
		a.callQuit()
	})
	flex.AddItem(mainMenu, 0, 1, true)
	flex.AddItem(help, 0, 3, false)
	a.callPage(flex)
}
