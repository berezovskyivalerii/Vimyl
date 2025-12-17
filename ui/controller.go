package ui

import (
	"image/color"
	"os"

	"gio.test/files"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type UI struct {
	// Themes and resources
	Theme *material.Theme

	// Widget's states
	PathInput widget.Editor
	ListState widget.List

	// Data
	CurrentFiles []os.FileInfo
	ErrMessage   string
}

// TODO
func NewUI(th *material.Theme) *UI {
	ui := &UI{Theme: th}
	ui.PathInput.SingleLine = true
	ui.PathInput.Submit = true
	ui.PathInput.SetText("/home/vasilisk/Music")

	ui.ListState.Axis = layout.Vertical
	return ui
}

// TODO
func (ui *UI) Update(gtx layout.Context) {
	// Check for events using Update(gtx)
	for {
		e, ok := ui.PathInput.Update(gtx)
		if !ok {
			break
		}
		if _, ok := e.(widget.SubmitEvent); ok {
			// User pressed Enter
			f, err := files.ListDir(ui.PathInput.Text())
			if err != nil {
				ui.ErrMessage = err.Error()
				ui.CurrentFiles = nil
			} else {
				ui.ErrMessage = ""
				ui.CurrentFiles = f
			}
		}
	}
}

// TODO
func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			// Header
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.H6(ui.Theme, "Music Directory:").Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

			// Input Field
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				ed := material.Editor(ui.Theme, &ui.PathInput, "/home/user/Music")
				border := widget.Border{
					Color:        color.NRGBA{A: 255},
					CornerRadius: unit.Dp(4),
					Width:        unit.Dp(1),
				}
				return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
				})
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

			// Error Message, if no error - shows nothing
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.ErrMessage != "" {
					lbl := material.Body2(ui.Theme, ui.ErrMessage)
					lbl.Color = color.NRGBA{R: 200, A: 255}
					return lbl.Layout(gtx)
				}
				return layout.Dimensions{}
			}),

			// File List
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				listStyle := material.List(ui.Theme, &ui.ListState)
				return listStyle.Layout(gtx, len(ui.CurrentFiles), func(gtx layout.Context, index int) layout.Dimensions {
					file := ui.CurrentFiles[index]
					txt := file.Name()
					if file.IsDir() {
						txt = "[DIR] " + txt
					}
					return layout.UniformInset(unit.Dp(4)).Layout(gtx,
						material.Body1(ui.Theme, txt).Layout,
					)
				})
			}),
		)
	})
}
