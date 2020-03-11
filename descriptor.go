// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bmfont

import (
	"image"
	"io"
	"os"
	"path/filepath"
)

// A Descriptor holds metadata for a bitmap font.
type Descriptor struct {
	Info    Info
	Common  Common
	Pages   map[int]Page
	Chars   map[rune]Char
	Kerning map[CharPair]Kerning
}

type Info struct {
	Face     string
	Size     int
	Bold     bool
	Italic   bool
	Charset  string
	Unicode  bool
	StretchH int
	Smooth   bool
	AA       int
	Padding  Padding
	Spacing  Spacing
	Outline  int
}

type Padding struct {
	Up, Right, Down, Left int
}

type Spacing struct {
	Horizontal, Vertical int
}

type Common struct {
	LineHeight   int
	Base         int
	ScaleW       int
	ScaleH       int
	Packed       bool
	AlphaChannel ChannelInfo
	RedChannel   ChannelInfo
	GreenChannel ChannelInfo
	BlueChannel  ChannelInfo
}

func (c *Common) Scale() image.Point {
	return image.Pt(c.ScaleH, c.ScaleH)
}

type ChannelInfo int

const (
	Glyph ChannelInfo = iota
	Outline
	GlyphAndOutline
	Zero
	One
)

type Page struct {
	ID   int
	File string
}

type Char struct {
	ID       rune
	X        int
	Y        int
	Width    int
	Height   int
	XOffset  int
	YOffset  int
	XAdvance int
	Page     int
	Channel  Channel
}

func (c *Char) Pos() image.Point {
	return image.Pt(c.X, c.Y)
}

func (c *Char) Size() image.Point {
	return image.Pt(c.Width, c.Height)
}

func (c *Char) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: c.Pos(),
		Max: c.Pos().Add(c.Size()),
	}
}

func (c *Char) Offset() image.Point {
	return image.Pt(c.XOffset, c.YOffset)
}

type Channel int

const (
	Blue  Channel = 1
	Green Channel = 2
	Red   Channel = 4
	Alpha Channel = 8
	All   Channel = 15
)

type CharPair struct {
	First, Second rune
}

type Kerning struct {
	Amount int
}

// LoadDescriptor loads the font descriptor data from a BMFont descriptor file in
// text format (usually with the file extension .fnt). It does not load the
// referenced page sheet images. If you also want to load the page sheet
// images, use the Load function to get a complete BitmapFont instance.
func LoadDescriptor(path string) (d *Descriptor, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closeChecked(f, &err)
	return parseDescriptor(filepath.Base(path), f)
}

// ReadDescriptor parses font descriptor data in BMFont's text format from a
// reader. It does not load the referenced page sheet images. If you also want
// to load the page sheet images, use the Load function to get a complete
// BitmapFont instance.
func ReadDescriptor(r io.Reader) (d *Descriptor, err error) {
	return parseDescriptor("bmfont", r)
}
