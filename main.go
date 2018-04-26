package main

import (
// "fmt"
)

func main() {
	image := MakeImage(500, 500)
	t := MakeMatrix(4, 4)
	t.Ident()
	e := MakeMatrix(4, 0)
	p := MakeMatrix(4, 0)
	// ParseFile("galleryscript", t, p, e, image)
	ParseFile("script", t, p, e, image)
}
