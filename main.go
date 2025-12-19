package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"log"
	"os"

	"gio.test/ui"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Gio Player"))
		w.Option(app.Size(600, 800))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	// Theme holds the general theme of app (fonts, icons, textSize, etc...)
	th := material.NewTheme()
	// Shaper converts strings of text into glyphs that can be displayed
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	myUI := ui.NewUI(th)

	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			myUI.Update(gtx)

			myUI.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
