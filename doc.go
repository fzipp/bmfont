// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bmfont reads bitmap fonts created with AngelCode's bitmap font
// generator or other tools that generate output in the same format, and
// draws texts with these fonts on images.
//
// The parser for the font descriptor files (.fnt) reads the text format, not the
// binary format. Format description:
// https://www.angelcode.com/products/bmfont/doc/file_format.html
package bmfont
