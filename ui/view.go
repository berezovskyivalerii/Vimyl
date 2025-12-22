package ui

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	ColorBackground = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
	ColorSurface    = color.NRGBA{R: 45, G: 45, B: 45, A: 255}
	ColorText       = color.NRGBA{R: 249, G: 247, B: 247, A: 255}
	ColorAccent     = color.NRGBA{R: 100, G: 149, B: 237, A: 255}
	ColorBorder     = color.NRGBA{R: 60, G: 60, B: 60, A: 255}
)

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, ColorBackground)

	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(ui.layoutHeader),
			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),

			layout.Flexed(1, ui.layoutSplitView),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),

			layout.Rigid(ui.layoutControls),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),

			layout.Rigid(ui.layoutInputZone),
		)
	})
}

func (ui *UI) layoutHeader(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			h6 := material.H6(ui.Theme, "Vimyl")
			h6.Color = ColorText
			h6.Font.Weight = 700
			return h6.Layout(gtx)
		}),
	)
}

func (ui *UI) layoutControls(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			pos, dur := ui.Audio.GetProgress()

			current := fmt.Sprintf("%d:%d", int64(pos/60), int64(pos)%60)
			duration := fmt.Sprintf("%d:%d", int64(dur/60), int64(dur)%60)

			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(ui.Theme, current)
					label.Color = ColorText
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(ui.Theme, duration)
					label.Color = ColorText
					return label.Layout(gtx)
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			s := material.Slider(ui.Theme, &ui.ProgressSlider)
			s.Color = ColorAccent
			return s.Layout(gtx)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			text := "Pause"
			if !ui.Audio.IsPlaying {
				text = "Resume"
			}

			btn := material.Button(ui.Theme, &ui.BtnPause, text)
			btn.Background = ColorAccent
			btn.Color = ColorBackground
			return btn.Layout(gtx)
		}),

		layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
	)
}

// Input Zone
func (ui *UI) layoutInputZone(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(ui.Theme, &ui.PathInput, "Path...")
			ed.Color = ColorText
			ed.HintColor = color.NRGBA{R: 100, G: 100, B: 100, A: 255}

			return widget.Border{
				Color:        ColorBorder,
				CornerRadius: unit.Dp(8),
				Width:        unit.Dp(1),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, ed.Layout)
			})
		}),

		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

		// Errors
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if ui.ErrMessage != "" {
				lbl := material.Body2(ui.Theme, ui.ErrMessage)
				lbl.Color = color.NRGBA{R: 255, G: 80, B: 80, A: 255}
				return lbl.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

// Split View (Folders | Files)
func (ui *UI) layoutSplitView(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		// Folders (25%)
		layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
			return drawList(gtx, ui.Theme, &ui.DListWrapper, &ui.DListState,
				len(ui.CurrentDirs), ui.DSelectedIndex, func(i int) string {
					return ui.CurrentDirs[i].Name()
				})
		}),

		// Files (75%)
		layout.Flexed(0.75, func(gtx layout.Context) layout.Dimensions {
			return drawList(gtx, ui.Theme, &ui.MListWrapper, &ui.MListState,
				len(ui.CurrentFiles), ui.MSelectedIndex, func(i int) string {
					return ui.CurrentFiles[i].Name()
				})
		}),
	)
}
