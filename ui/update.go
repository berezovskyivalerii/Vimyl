package ui

import (
	"log/slog"
	"path/filepath"

	"gio.test/files"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
)

func (ui *UI) Update(gtx layout.Context) {
	prevDirIndex := ui.DSelectedIndex

	// Logic for folders list
	ui.handleListNav(gtx, &ui.DListWrapper, &ui.DSelectedIndex, len(ui.CurrentDirs), &ui.DListState, "DIRS")

	shouldUpdatePrewiew := prevDirIndex != ui.DSelectedIndex

	if len(ui.CurrentFiles) == 0 && len(ui.CurrentDirs) > 0 {
		shouldUpdatePrewiew = true
	}

	if shouldUpdatePrewiew && len(ui.CurrentDirs) > 0 {
		seletedDir := ui.CurrentDirs[ui.DSelectedIndex]

		// Check for General folder
		if seletedDir.Name() == "General" { // If true than search by default path
			generalFiles, err := files.ListFiles(ui.PathInput.Text())
			if err == nil {
				ui.CurrentFiles = generalFiles
				ui.MSelectedIndex = 0
				ui.MListState.ScrollTo(0)
			} else {
				ui.CurrentFiles = nil
			}
		} else { // In other way search by specific one
			fullPath := filepath.Join(ui.PathInput.Text(), seletedDir.Name())
			newFiles, err := files.ListFiles(fullPath)
			if err == nil {
				ui.CurrentFiles = newFiles
				ui.MSelectedIndex = 0
				ui.MListState.ScrollTo(0)
			} else {
				ui.CurrentFiles = nil
			}
		}
	}

	if ui.ProgressSlider.Update(gtx) {
		ui.Audio.Seek(float64(ui.ProgressSlider.Value))
	}

	if ui.Audio.IsPlaying && !ui.ProgressSlider.Dragging() {
		pos, dur := ui.Audio.GetProgress()
		if dur > 0 {
			ui.ProgressSlider.Value = float32(pos / dur)
			// Force interface redraw to animate slider smoothly
			gtx.Execute(op.InvalidateCmd{})
		}
	}
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

	for ui.BtnPause.Clicked(gtx) {
		ui.Audio.TogglePause()
		slog.Info("Button has been pressed")
	}
	ui.handlePathInput(gtx)
}

// Helper to do not duplicate logic for each list
func (ui *UI) handleListNav(gtx layout.Context, tag *widget.Clickable, index *int, length int, listState *widget.List, logName string) {
	for {
		ev, ok := gtx.Event(key.Filter{Focus: tag})
		if !ok {
			break
		}
		if ke, ok := ev.(key.Event); ok && ke.State == key.Press {
			switch ke.Name {
			case key.NameLeftArrow, "H":
				gtx.Execute(key.FocusCmd{Tag: &ui.DListWrapper})
			case key.NameRightArrow, "L":
				gtx.Execute(key.FocusCmd{Tag: &ui.MListWrapper})
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
			case key.NameSpace:
				ui.Audio.TogglePause()
			case key.NameReturn, "Enter":
				if len(ui.CurrentFiles) > 0 {
					f := ui.CurrentFiles[ui.MSelectedIndex]
					fullPath := filepath.Join(ui.PathInput.Text(), f.Name())
					ui.Audio.PlayFile(fullPath)
				}
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
				ui.CurrentDirs = dirsList
				ui.DSelectedIndex = 0

				if len(dirsList) > 0 {
					firstDir := dirsList[0]
					previewPath := filepath.Join(path, firstDir.Name())
					previewFiles, _ := files.ListFiles(previewPath)
					ui.CurrentFiles = previewFiles
				} else {
					ui.CurrentFiles = filesList
				}

				ui.MSelectedIndex = 0
				// Focus to folders by default
				gtx.Execute(key.FocusCmd{Tag: &ui.DListWrapper})
			}
		}
	}
}
