// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/system8bit/ddzipper/dropfiles/zipper"
	"strings"
)

func main() {
	var textEdit *walk.TextEdit
	MainWindow{
		Title:   "Folder Zipper",
		MinSize: Size{320, 240},
		Size:    Size{320, 240},
		Layout:  VBox{},
		OnDropFiles: func(files []string) {
			textEdit.SetText(strings.Join(files, "\r\n"))
			for _, file := range files {
				err := zipper.CompressFolder(file)
				if err != nil {
					return
				}
			}
		},
		Children: []Widget{
			TextEdit{
				AssignTo: &textEdit,
				ReadOnly: true,
				Text:     "Drop files here, from windows explorer...",
			},
		},
	}.Run()
}
