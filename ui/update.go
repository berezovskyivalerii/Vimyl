package ui

import (
	"log/slog"

	"gio.test/files"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"
)

func (ui *UI) Update(gtx layout.Context) {
	// Logic for folders list
	ui.handleListNav(gtx, &ui.DListWrapper, &ui.DSelectedIndex, len(ui.CurrentDirs), &ui.DListState, "DIRS")
	// Logic for files list
	ui.handleListNav(gtx, &ui.MListWrapper, &ui.MSelectedIndex, len(ui.CurrentFiles), &ui.MListState, "FILES")

	// Move focus to files
	for ui.MListWrapper.Clicked(gtx) {
		gtx.Execute(key.FocusCmd{Tag: &ui.MListWrapper})
	}
	// Move focus to folders
	for ui.DListWrapper.Clicked(gtx) {
		gtx.Execute(key.FocusCmd{Tag: &ui.DListWrapper})
	}

	ui.handlePathInput(gtx)
}

// Helper to do not duplicate lofic for each list
func (ui *UI) handleListNav(gtx layout.Context, tag *widget.Clickable, index *int, length int, listState *widget.List, logName string) {
	for {
		ev, ok := gtx.Event(key.Filter{Focus: tag})
		if !ok {
			break
		}
		if ke, ok := ev.(key.Event); ok && ke.State == key.Press {
			switch ke.Name {
			case key.NameDownArrow, "J":
				if *index < length-1 {
					*index++
					listState.ScrollTo(*index)
				}
			case key.NameUpArrow, "K":
				if *index > 0 {
					*index--
					listState.ScrollTo(*index)
				}
			case key.NameHome:
				*index = 0
				listState.ScrollTo(0)
			case key.NameEnd:
				*index = length - 1
				listState.ScrollTo(*index)
			}
			slog.Info("Nav Key Pressed", "list", logName, "key", ke.Name)
		}
	}
}

// Handle path input
func (ui *UI) handlePathInput(gtx layout.Context) {
	for {
		e, ok := ui.PathInput.Update(gtx)
		if !ok {
			break
		}
		if _, ok := e.(widget.SubmitEvent); ok {
			path := ui.PathInput.Text()

			filesList, err := files.ListFiles(path)
			dirsList, err2 := files.ListDirs(path)

			if err != nil {
				ui.ErrMessage = err.Error()
				ui.CurrentFiles = nil
			} else if err2 != nil {
				ui.ErrMessage = err2.Error()
				ui.CurrentDirs = nil
			} else {
				ui.ErrMessage = ""
				ui.CurrentFiles = filesList
				ui.CurrentDirs = dirsList

				ui.MSelectedIndex = 0
				ui.DSelectedIndex = 0

				// Focus to folders by default
				gtx.Execute(key.FocusCmd{Tag: &ui.DListWrapper})
			}
		}
	}
}
