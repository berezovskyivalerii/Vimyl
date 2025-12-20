package ui

import (
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
	// Refister Tag for events
	event.Op(gtx.Ops, wrapper)

	// Check focus
	isFocused := gtx.Source.Focused(wrapper)
	borderColor := color.NRGBA{A: 0}
	if isFocused {
		borderColor = color.NRGBA{R: 255, A: 255}
	}

	// Wrapper for list
	return wrapper.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return widget.Border{
			Color:        borderColor,
			CornerRadius: unit.Dp(0),
			Width:        unit.Dp(2),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// List by itself
			listStyle := material.List(th, listState)
			return listStyle.Layout(gtx, count, func(gtx layout.Context, index int) layout.Dimensions {
				// Here we get text from helper func
				txt := namer(index)
				return layout.Stack{}.Layout(gtx,
					// Color background when selected
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						if index == selectedIndex && isFocused {
							defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
							paint.Fill(gtx.Ops, color.NRGBA{R: 220, G: 220, B: 220, A: 255})
						}
						return layout.Dimensions{Size: gtx.Constraints.Min}
					}),

					// Text
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(5)).Layout(gtx,
							material.Body1(th, txt).Layout,
						)
					}),
				)
			})
		})
	})
}
