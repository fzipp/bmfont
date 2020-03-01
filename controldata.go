// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bmfont implements a parser for bitmap font control files (.fnt)
// created with AngelCode's bitmap font generator or other tools that generate
// output in the same format.
//
// The parser parses the text format, not the binary format. Format description:
// http://www.angelcode.com/products/bmfont/doc/file_format.html
package bmfont

import "image"

type ControlData struct {
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
	First  rune
	Second rune
}

type Kerning struct {
	Amount int
}
