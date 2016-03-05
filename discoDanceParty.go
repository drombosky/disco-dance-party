// Generate documentation.
//go:generate go get github.com/robertkrimen/godocdown/godocdown
//go:generate godocdown -output=README.md
//go:generate godocdown -output=hue/README.md hue
//go:generate godocdown -output=hue/light/README.md hue/light

// Disco Dance Party is program used for controlling Hue lights and run animations.
package main

func main() {}
