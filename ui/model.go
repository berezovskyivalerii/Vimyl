package ui

import (
	"os"

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

	ErrMessage string
}

func NewUI(th *material.Theme) *UI {
	ui := &UI{Theme: th}

	ui.PathInput.SingleLine = true
	ui.PathInput.Submit = true
	ui.PathInput.SetText("/home")

	ui.DListState.Axis = layout.Vertical
	ui.MListState.Axis = layout.Vertical

	return ui
}
