package ui

import (
	"strconv"

	"github.com/rivo/tview"
)

func (a *KirinukiApp) callPage(content tview.Primitive) {
	f := tview.NewFrame(content)
	f.SetBorder(true)
	f.SetBorderColor(t_colors[frame_border_color])
	f.SetBorders(2, 2, 1, 1, 1, 1)
	f.AddText("Kirinuki -  @LosAngeles971", true, tview.AlignLeft, t_colors[frame_header_color])
	f.AddText("Press (ESC) key to have the main menù", false, tview.AlignCenter, t_colors[frame_footer_color])
	a.app.SetRoot(f, true)
}

func getInteger(f *tview.InputField) int {
	a, err := strconv.Atoi(f.GetText())
	if err != nil {
		return 0
	}
	return a
}

func getFloat(f *tview.InputField) float64 {
	a, err := strconv.ParseFloat(f.GetText(), 64)
	if err != nil {
		return 0.0
	}
	return a
}

