package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Top bar with path input
			layout.Rigid(ui.layoutTopBar),
			// Folders and files
			layout.Flexed(1, ui.layoutSplitView),
		)
	})
}

// Input
func (ui *UI) layoutTopBar(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.H6(ui.Theme, "File Browser").Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(ui.Theme, &ui.PathInput, "Path...")
			return widget.Border{
				Color:        color.NRGBA{A: 200},
				CornerRadius: unit.Dp(4),
				Width:        unit.Dp(1),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if ui.ErrMessage != "" {
				lbl := material.Body2(ui.Theme, ui.ErrMessage)
				lbl.Color = color.NRGBA{R: 200, A: 255}
				return lbl.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

// Folders and Files
func (ui *UI) layoutSplitView(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		// Folders (25% width)
		layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
			return drawList(gtx, ui.Theme, &ui.DListWrapper, &ui.DListState,
				len(ui.CurrentDirs), ui.DSelectedIndex, func(i int) string {
					return "[[DIR]]\t" + ui.CurrentDirs[i].Name()
				})
		}),

		// Files (75% width)
		layout.Flexed(0.75, func(gtx layout.Context) layout.Dimensions {
			return drawList(gtx, ui.Theme, &ui.MListWrapper, &ui.MListState,
				len(ui.CurrentFiles), ui.MSelectedIndex, func(i int) string {
					return "[[MUS]]\t" + ui.CurrentFiles[i].Name()
				})
		}),
	)
}
