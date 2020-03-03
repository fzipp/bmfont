// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bmfont

import (
	"image"
	"image/draw"
	"os"
	"path/filepath"
)

type BitmapFont struct {
	Ctl    *ControlData
	Sheets map[int]image.Image
}

func Load(path string) (*BitmapFont, error) {
	ctl, err := LoadControlData(path)
	if err != nil {
		return nil, err
	}
	font := BitmapFont{
		Ctl:    ctl,
		Sheets: make(map[int]image.Image),
	}
	for i, page := range ctl.Pages {
		dir, _ := filepath.Split(path)
		sheetPath := filepath.Join(dir, page.File)
		sheetFile, err := os.Open(sheetPath)
		if err != nil {
			return nil, err
		}
		sheet, _, err := image.Decode(sheetFile)
		if err != nil {
			sheetFile.Close()
			return nil, err
		}
		sheetFile.Close()
		font.Sheets[i] = sheet
	}
	return &font, nil
}

func (f *BitmapFont) DrawText(dst draw.Image, pos image.Point, text string) {
	cursor := pos
	var prev rune
	for i, r := range text {
		if r == '\n' {
			cursor.X = pos.X
			cursor.Y += f.Ctl.Common.LineHeight
			continue
		}
		ch, ok := f.char(r)
		if !ok {
			continue
		}
		sheet := f.Sheets[ch.Page]
		min := image.Pt(
			cursor.X+ch.XOffset,
			cursor.Y-f.Ctl.Common.Base+ch.YOffset,
		)
		dr := image.Rectangle{
			Min: min,
			Max: min.Add(ch.Size()),
		}
		draw.Draw(dst, dr, sheet, ch.Pos(), draw.Over)

		cursor.X += ch.XAdvance
		if i > 0 {
			pair := CharPair{First: prev, Second: r}
			kerning, ok := f.Ctl.Kerning[pair]
			if ok {
				cursor.X += kerning.Amount
			}
		}
		prev = r
	}
}

func (f *BitmapFont) char(r rune) (c Char, ok bool) {
	c, ok = f.Ctl.Chars[r]
	if !ok {
		c, ok = f.Ctl.Chars['?']
		return c, ok
	}
	return c, ok
}
