// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bmfont

import (
	"image"
	"image/draw"
	"io"
	"os"
	"path/filepath"
)

type BitmapFont struct {
	Descriptor *Descriptor
	PageSheets map[int]image.Image
}

// Load loads a bitmap font from a BMFont descriptor file (.fnt) in text format
// including all the referenced page sheet images. The resulting bitmap font
// is ready to be used to draw text on an image.
func Load(path string) (f *BitmapFont, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closeChecked(file, &err)
	dir, _ := filepath.Split(path)
	return Read(file, fileSheets(dir))
}

// Read reads a bitmap font from a BMFont descriptor in text format including
// all the referenced page sheet images.
// The page sheet images are read from the readers provided by the given
// SheetReaderFunc. These sheet readers are closed after use. If you want to
// keep them open wrap them via ioutil.NopCloser.
// The resulting bitmap font is ready to be used to draw text on an image.
func Read(r io.Reader, sheets SheetReaderFunc) (f *BitmapFont, err error) {
	desc, err := ReadDescriptor(r)
	if err != nil {
		return nil, err
	}
	font := BitmapFont{
		Descriptor: desc,
		PageSheets: make(map[int]image.Image),
	}
	for i, page := range desc.Pages {
		sheet, err := sheets.read(page.File)
		if err != nil {
			return nil, err
		}
		font.PageSheets[i] = sheet
	}
	return &font, nil
}

type SheetReaderFunc func(filename string) (io.ReadCloser, error)

func (f SheetReaderFunc) read(filename string) (img image.Image, err error) {
	r, err := f(filename)
	if err != nil {
		return nil, err
	}
	defer closeChecked(r, &err)
	sheet, _, err := image.Decode(r)
	return sheet, err
}

func fileSheets(directory string) SheetReaderFunc {
	return func(filename string) (io.ReadCloser, error) {
		path := filepath.Join(directory, filename)
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

// DrawText draws the given text on the destination image starting at the
// given position. The start position is on the base line of the first line of
// text, and the characters usually extend above the base line.
// The text may contain newlines. Text with multiple lines is drawn left
// aligned.
func (f *BitmapFont) DrawText(dst draw.Image, pos image.Point, text string) {
	f.drawText(imageDrawer{dst}, pos, text)
}

// MeasureText calculates the bounding box for the given text as if it was
// drawn at position (0, 0). The Min point usually has a negative Y coordinate,
// since the start position of DrawText is on the base line and the characters
// extend above the base line. The X coordinate can also be negative, depending
// on the character offsets.
func (f *BitmapFont) MeasureText(text string) image.Rectangle {
	var m boundsMeasurer
	f.drawText(&m, image.Point{}, text)
	return m.bounds
}

func (f *BitmapFont) drawText(dst drawer, pos image.Point, text string) {
	cursor := pos
	var prev rune
	for i, r := range text {
		if r == '\n' {
			cursor.X = pos.X
			cursor.Y += f.Descriptor.Common.LineHeight
			continue
		}
		ch, ok := f.char(r)
		if !ok {
			continue
		}
		sheet := f.PageSheets[ch.Page]
		min := image.Pt(
			cursor.X+ch.XOffset,
			cursor.Y-f.Descriptor.Common.Base+ch.YOffset,
		)
		dr := image.Rectangle{
			Min: min,
			Max: min.Add(ch.Size()),
		}
		dst.Draw(dr, sheet, ch.Pos())

		cursor.X += ch.XAdvance
		if i > 0 {
			pair := CharPair{First: prev, Second: r}
			kerning, ok := f.Descriptor.Kerning[pair]
			if ok {
				cursor.X += kerning.Amount
			}
		}
		prev = r
	}
}

func (f *BitmapFont) char(r rune) (c Char, ok bool) {
	c, ok = f.Descriptor.Chars[r]
	if !ok {
		c, ok = f.Descriptor.Chars['?']
		return c, ok
	}
	return c, ok
}

type drawer interface {
	Draw(r image.Rectangle, src image.Image, sp image.Point)
}

type imageDrawer struct {
	draw.Image
}

func (dst imageDrawer) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	draw.Draw(dst, r, src, sp, draw.Over)
}

type boundsMeasurer struct {
	bounds image.Rectangle
}

func (m *boundsMeasurer) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	_, _ = src, sp
	m.bounds = m.bounds.Union(r)
}

func closeChecked(c io.Closer, err *error) {
	cErr := c.Close()
	if cErr != nil && *err == nil {
		*err = cErr
	}
}
