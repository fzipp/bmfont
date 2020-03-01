// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bmfont

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"
)

// LoadControlData loads the font control data from a file.
func LoadControlData(path string) (*ControlData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseControlData(filepath.Base(path), file)
}

func ReadControlData(r io.Reader) (*ControlData, error) {
	return parseControlData("bmfont", r)
}

func parseControlData(filename string, r io.Reader) (*ControlData, error) {
	var p tagsParser
	tags, err := p.parse(filename, r)
	if err != nil {
		return nil, err
	}
	font := ControlData{
		Pages:   make(map[int]Page),
		Chars:   make(map[rune]Char),
		Kerning: make(map[CharPair]Kerning),
	}
	for _, tag := range tags {
		switch tag.name {
		case "info":
			font.Info = Info{
				Face:     tag.stringAttr("face"),
				Size:     tag.intAttr("size"),
				Bold:     tag.boolAttr("bold"),
				Italic:   tag.boolAttr("italic"),
				Charset:  tag.stringAttr("charset"),
				Unicode:  tag.boolAttr("unicode"),
				StretchH: tag.intAttr("stretchH"),
				Smooth:   tag.boolAttr("smooth"),
				AA:       tag.intAttr("aa"),
				Padding:  paddingFrom(tag.intListAttr("padding", 4)),
				Spacing:  spacingFrom(tag.intListAttr("spacing", 2)),
				Outline:  tag.intAttr("outline"),
			}
		case "common":
			font.Common = Common{
				LineHeight:   tag.intAttr("lineHeight"),
				Base:         tag.intAttr("base"),
				ScaleW:       tag.intAttr("scaleW"),
				ScaleH:       tag.intAttr("scaleH"),
				Packed:       tag.boolAttr("packed"),
				AlphaChannel: ChannelInfo(tag.intAttr("alphaChnl")),
				RedChannel:   ChannelInfo(tag.intAttr("redChnl")),
				GreenChannel: ChannelInfo(tag.intAttr("greenChnl")),
				BlueChannel:  ChannelInfo(tag.intAttr("blueChnl")),
			}
		case "page":
			id := tag.intAttr("id")
			font.Pages[id] = Page{
				ID:   id,
				File: tag.stringAttr("file"),
			}
		case "char":
			id := rune(tag.intAttr("id"))
			font.Chars[id] = Char{
				ID:       id,
				X:        tag.intAttr("x"),
				Y:        tag.intAttr("y"),
				Width:    tag.intAttr("width"),
				Height:   tag.intAttr("height"),
				XOffset:  tag.intAttr("xoffset"),
				YOffset:  tag.intAttr("yoffset"),
				XAdvance: tag.intAttr("xadvance"),
			}
		case "kerning":
			pair := CharPair{
				First:  rune(tag.intAttr("first")),
				Second: rune(tag.intAttr("second")),
			}
			font.Kerning[pair] = Kerning{
				Amount: tag.intAttr("amount"),
			}
		}
	}
	return &font, nil
}

func paddingFrom(values []int) Padding {
	return Padding{
		Up:    values[0],
		Right: values[1],
		Down:  values[2],
		Left:  values[3],
	}
}

func spacingFrom(values []int) Spacing {
	return Spacing{
		Horizontal: values[0],
		Vertical:   values[1],
	}
}

type tagsParser struct {
	errors  errorList
	scanner scanner.Scanner
	pos     scanner.Position
	tok     rune
	lit     string
}

func (p *tagsParser) next() {
	p.tok = p.scanner.Scan()
	p.pos = p.scanner.Position
	p.lit = p.scanner.TokenText()
}

func (p *tagsParser) parse(filename string, r io.Reader) ([]tag, error) {
	p.scanner.Init(r)
	p.scanner.Filename = filename
	p.scanner.Whitespace ^= 1 << '\n'
	p.scanner.Error = func(s *scanner.Scanner, msg string) {}
	p.next()

	var tags []tag
	for p.tok != scanner.EOF {
		tagName := p.lit
		p.expect(scanner.Ident, "tag name")
		var attrs = make(map[string]string)
		for p.tok != '\n' && p.tok != scanner.EOF {
			attrName := p.lit
			p.expect(scanner.Ident, "attribute name")
			value := ""
			var err error
			p.expect('=', `"="`)
			switch p.tok {
			case scanner.String:
				value, err = strconv.Unquote(p.lit)
				if err != nil {
					// end this line, rest is garbage
					p.tok = '\n'
					continue
				}
				if p.scanner.Peek() == '"' {
					// workaround for `letter="""`
					p.scanner.Next()
					value += `"`
				}
				p.next()
			case scanner.Int, '-':
				value = p.parseIntList()
			default:
				p.errorExpected("string or integer attribute value")
			}
			attrs[attrName] = value
		}
		tags = append(tags, tag{
			name:  tagName,
			attrs: attrs,
		})
		p.next()
	}
	return tags, p.errors.Err()
}

func (p *tagsParser) parseIntList() string {
	var sb strings.Builder
	for p.tok == scanner.Int || p.tok == ',' || p.tok == '-' {
		sb.WriteString(p.lit)
		p.next()
	}
	return sb.String()
}

func (p *tagsParser) expect(tok rune, msg string) {
	if p.tok != tok {
		p.errorExpected(msg)
	}
	p.next()
}

func (p *tagsParser) errorExpected(msg string) {
	p.error(newError(p.pos, "expected "+msg+", found "+scanner.TokenString(p.tok)))
}

func (p *tagsParser) error(err error) {
	p.errors = append(p.errors, err)
}

func newError(pos scanner.Position, msg string) error {
	return errors.New(pos.String() + ": " + msg)
}

type tag struct {
	name  string
	attrs map[string]string
}

func (t *tag) intAttr(name string) int {
	value, _ := strconv.Atoi(t.stringAttr(name))
	return value
}

func (t *tag) stringAttr(name string) string {
	return t.attrs[name]
}

func (t *tag) boolAttr(name string) bool {
	return t.intAttr(name) != 0
}

func (t *tag) intListAttr(name string, n int) []int {
	values := make([]int, n)
	parts := strings.Split(t.stringAttr(name), ",")
	for i, part := range parts {
		if i == len(values) {
			break
		}
		value, _ := strconv.Atoi(strings.TrimSpace(part))
		values[i] = value
	}
	return values
}

type errorList []error

func (list errorList) Err() error {
	if len(list) == 0 {
		return nil
	}
	return list
}

func (list errorList) Error() string {
	if len(list) == 0 {
		return "no errors"
	}
	return list[0].Error()
}
