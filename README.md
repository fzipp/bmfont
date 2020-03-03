# bmfont

A Go package to load and render bitmap fonts
created with [AngelCode's bitmap font generator](https://www.angelcode.com/products/bmfont/)
or other tools that generate output in the same format.

This package uses the [text format](ttp://www.angelcode.com/products/bmfont/doc/file_format.html)
for font control files (.fnt), not the binary format.

## Documentation

Package documentation is available [on pkg.go.dev](https://pkg.go.dev/github.com/fzipp/bmfont?tab=doc).

## Example usage

Load a bitmap font and draw text to an image:

```
package main

import (
	"fmt"
	"log"

	"github.com/fzipp/bmfont"
)

func main() {
	font, err := bmfont.Load("ExampleFont.fnt")
	if err != nil {
		log.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 600, 300))
	font.DrawText(img, image.Pt(10, 20), `hello, world
This is an example.
abcdefghijklmnopqrstuvwxyz
ABCDEFGHIJKLMNOPQRSTUVWXYZ`)
	// ...
}
```

Only load the control data of a bitmap font:

```
package main

import (
	"fmt"
	"log"

	"github.com/fzipp/bmfont"
)

func main() {
	font, err := bmfont.LoadControlData("ExampleFont.fnt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(font.Info.Face)
	fmt.Println("line height:", font.Common.LineHeight)
	fmt.Println("letter A width:", font.Chars['A'].Width)
}
```

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
