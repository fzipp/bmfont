# bmfont

Implements a parser for bitmap font control files (.fnt)
created with [AngelCode's bitmap font generator](https://www.angelcode.com/products/bmfont/)
or other tools that generate output in the same format.

The parser parses the [text format](ttp://www.angelcode.com/products/bmfont/doc/file_format.html), not the binary format.

## Documentation

Package documentation is available [on pkg.go.dev](https://pkg.go.dev/github.com/fzipp/bmfont?tab=doc).

## Example usage

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
