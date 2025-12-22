package ui

import (
	"image"
	"image/color"

	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func drawList(gtx layout.Context, th *material.Theme, wrapper *widget.Clickable, listState *widget.List, count int, selectedIndex int, namer func(int) string) layout.Dimensions {
	event.Op(gtx.Ops, wrapper)

	isFocused := gtx.Source.Focused(wrapper)

	borderColor := ColorBorder
	if isFocused {
		borderColor = ColorAccent
	}

	// Wrapper
	return wrapper.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{
			Color:        borderColor,
			CornerRadius: unit.Dp(8),
			Width:        unit.Dp(1),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				listStyle := material.List(th, listState)
				listStyle.AnchorStrategy = material.Overlay

				return listStyle.Layout(gtx, count, func(gtx layout.Context, index int) layout.Dimensions {
					txt := namer(index)

					textColor := ColorText
					bgColor := color.NRGBA{A: 0}

					if index == selectedIndex {
						if isFocused {
							bgColor = ColorAccent
							textColor = ColorBackground
						} else {
							bgColor = ColorSurface
						}
					}

					return layout.Stack{}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 4).Push(gtx.Ops).Pop()
							paint.Fill(gtx.Ops, bgColor)
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}),

						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body1(th, txt)
								lbl.Color = textColor
								return lbl.Layout(gtx)
							})
						}),
					)
				})
			})
		})
	})
}

