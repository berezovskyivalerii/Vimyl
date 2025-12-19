package ui

import (
	"image/color"
	"log/slog"
	"os"

	"gio.test/files"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type UI struct {
	Theme *material.Theme

	PathInput   widget.Editor
	ListState   widget.List
	ListWrapper widget.Clickable

	CurrentFiles []os.FileInfo
	ErrMessage   string

	SelectedIndex int
}

func NewUI(th *material.Theme) *UI {
	ui := &UI{Theme: th}
	ui.PathInput.SingleLine = true
	ui.PathInput.Submit = true
	ui.PathInput.SetText("/home")

	ui.ListState.Axis = layout.Vertical
	return ui
}

func (ui *UI) Update(gtx layout.Context) {
	// 1. List
	for {
		ev, ok := gtx.Event(key.Filter{
			Focus: &ui.ListWrapper,
		})
		if !ok {
			break
		}

		// Check if this a key event or not
		if ke, ok := ev.(key.Event); ok {
			if ke.State == key.Press {
				// Navigation logic
				switch ke.Name {
				case key.NameDownArrow, "J": // ArrowDown or 'j'
					if ui.SelectedIndex < len(ui.CurrentFiles)-1 {
						ui.SelectedIndex++
						ui.ListState.ScrollTo(ui.SelectedIndex)
					}
					slog.Info("DownArrow or J PRESSED")
				case key.NameUpArrow, "K": // ArrowUp or 'k'
					if ui.SelectedIndex > 0 {
						ui.SelectedIndex--
						ui.ListState.ScrollTo(ui.SelectedIndex)
					}
					slog.Info("UpArrow or K PRESSED")
				case key.NameHome: // Home
					ui.SelectedIndex = 0
					ui.ListState.ScrollTo(0)

					slog.Info("Home PRESSED")
				case key.NameEnd: // End
					ui.SelectedIndex = len(ui.CurrentFiles) - 1
					ui.ListState.ScrollTo(ui.SelectedIndex)

					slog.Info("End PRESSED")
				}
			}
		}
	}

	// 2. Mouse click (focus back)
	for ui.ListWrapper.Clicked(gtx) {
		gtx.Execute(key.FocusCmd{Tag: &ui.ListWrapper})
	}

	// 3. Path Input
	for {
		e, ok := ui.PathInput.Update(gtx)
		if !ok {
			break
		}
		if _, ok := e.(widget.SubmitEvent); ok {
			f, err := files.ListDir(ui.PathInput.Text())
			if err != nil {
				ui.ErrMessage = err.Error()
				ui.CurrentFiles = nil
			} else {
				ui.ErrMessage = ""
				ui.CurrentFiles = f
				ui.SelectedIndex = 0

				// Move focus to list
				gtx.Execute(key.FocusCmd{Tag: &ui.ListWrapper})
			}
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Header
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.H6(ui.Theme, "File Browser").Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

			// Input
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

			// Error
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.ErrMessage != "" {
					lbl := material.Body2(ui.Theme, ui.ErrMessage)
					lbl.Color = color.NRGBA{R: 200, A: 255}
					return lbl.Layout(gtx)
				}
				return layout.Dimensions{}
			}),

			// Files list
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				// 1. Tag registration
				event.Op(gtx.Ops, &ui.ListWrapper)

				// 2. Focus visualization
				isFocused := gtx.Source.Focused(&ui.ListWrapper)
				borderColor := color.NRGBA{A: 0}
				if isFocused {
					borderColor = color.NRGBA{B: 255, A: 255}
				}

				return ui.ListWrapper.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return widget.Border{
						Color:        borderColor,
						CornerRadius: unit.Dp(0),
						Width:        unit.Dp(2),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						listStyle := material.List(ui.Theme, &ui.ListState)
						return listStyle.Layout(gtx, len(ui.CurrentFiles), func(gtx layout.Context, index int) layout.Dimensions {
							f := ui.CurrentFiles[index]

							return layout.Stack{}.Layout(gtx,
								// Background
								layout.Expanded(func(gtx layout.Context) layout.Dimensions {
									if index == ui.SelectedIndex && isFocused {
										defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
										paint.Fill(gtx.Ops, color.NRGBA{R: 220, G: 220, B: 220, A: 255})
									}
									return layout.Dimensions{Size: gtx.Constraints.Min}
								}),
								// Text
								layout.Stacked(func(gtx layout.Context) layout.Dimensions {
									name := f.Name()
									if f.IsDir() {
										name = "📂 " + name
									} else {
										name = "📄 " + name
									}
									return layout.UniformInset(unit.Dp(5)).Layout(gtx,
										material.Body1(ui.Theme, name).Layout,
									)
								}),
							)
						})
					})
				})
			}),
		)
	})
}
