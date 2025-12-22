package ui

import (
	"os"

	"gio.test/player"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type UI struct {
	Theme *material.Theme

	PathInput widget.Editor

	// Dirs panel (left)
	DListState     widget.List
	DListWrapper   widget.Clickable
	CurrentDirs    []os.FileInfo
	DSelectedIndex int

	// Files panel (right)
	MListState     widget.List
	MListWrapper   widget.Clickable
	CurrentFiles   []os.FileInfo
	MSelectedIndex int

	BtnPause widget.Clickable

	ErrMessage     string
	Audio          *player.AudioPlayer
	ProgressSlider widget.Float
}

func NewUI(th *material.Theme) *UI {
	ui := &UI{Theme: th}

	ui.PathInput.SingleLine = true
	ui.PathInput.Submit = true
	ui.PathInput.SetText("/home/vasilisk/Music")

	ui.DListState.Axis = layout.Vertical
	ui.MListState.Axis = layout.Vertical

	ui.Audio = player.NewAudioPlayer()
	return ui
}
