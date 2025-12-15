package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Helper function to read directory
func showDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entries, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("File Browser"),
			app.Size(unit.Dp(400), unit.Dp(600)), // Made window taller for the list
		)

		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	// 1. STATE VARIABLES
	var myInput widget.Editor
	myInput.SingleLine = true
	myInput.Submit = true

	var ops op.Ops

	// Data to display
	var files []os.FileInfo
	var errMessage string

	// List state
	var fileListState widget.List
	fileListState.Axis = layout.Vertical

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// --- 2. LOGIC (FIXED) ---
			// Check for events using Update(gtx)
			for {
				e, ok := myInput.Update(gtx)
				if !ok {
					break
				}
				if _, ok := e.(widget.SubmitEvent); ok {
					// User pressed Enter
					f, err := showDir(myInput.Text())
					if err != nil {
						errMessage = err.Error()
						files = nil
					} else {
						errMessage = ""
						files = f
					}
				}
			}

			// --- 3. DRAWING ---
			layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// Header
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.H6(th, "Music Directory:").Layout(gtx)
					}),

					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

					// Input Field
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						ed := material.Editor(th, &myInput, "/home/user/music")
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

					// Error Message
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if errMessage != "" {
							lbl := material.Body2(th, errMessage)
							lbl.Color = color.NRGBA{R: 200, A: 255}
							return lbl.Layout(gtx)
						}
						return layout.Dimensions{}
					}),

					// File List
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						listStyle := material.List(th, &fileListState)
						return listStyle.Layout(gtx, len(files), func(gtx layout.Context, index int) layout.Dimensions {
							file := files[index]
							txt := file.Name()
							if file.IsDir() {
								txt = "[DIR] " + txt
							}
							return layout.UniformInset(unit.Dp(4)).Layout(gtx,
								material.Body1(th, txt).Layout,
							)
						})
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}
