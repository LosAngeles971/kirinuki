package ui

import "github.com/gdamore/tcell/v2"

const (
	frame_border_color   = "frame_border_color"
	frame_header_color   = "frame_header_color"
	frame_footer_color   = "frame_footer_color"
	table_header_color   = "table_header_color"
	table_cell_color     = "table_cell_color"
	table_highcell_color = "table_highcell_color"
)

var t_colors = map[string]tcell.Color{
	frame_border_color:   tcell.ColorBlueViolet,
	frame_header_color:   tcell.ColorBlue,
	frame_footer_color:   tcell.ColorDarkRed,
	table_header_color:   tcell.ColorGreenYellow,
	table_cell_color:     tcell.ColorYellow,
	table_highcell_color: tcell.ColorLightCoral,
}
