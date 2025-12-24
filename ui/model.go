package ui

import (
	"os"

	"gio.test/player"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
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

	// Icons controls
	BtnPlay   widget.Clickable
	IconPlay  *widget.Icon
	IconPause *widget.Icon

	BtnPrev  widget.Clickable
	IconPrev *widget.Icon

	BtnNext  widget.Clickable
	IconNext *widget.Icon

	// Error
	ErrMessage string

	// Audio settings
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

	ui.IconPlay, _ = widget.NewIcon(icons.AVPlayArrow)
	ui.IconPause, _ = widget.NewIcon(icons.AVPause)
	ui.IconPrev, _ = widget.NewIcon(icons.AVSkipPrevious)
	ui.IconNext, _ = widget.NewIcon(icons.AVSkipNext)

	ui.Audio = player.NewAudioPlayer()
	return ui
}
