// Generate documentation.
//go:generate go get github.com/robertkrimen/godocdown/godocdown
//go:generate godocdown -output=README.md
//go:generate godocdown -output=hue/README.md hue
//go:generate godocdown -output=hue/client/README.md hue/client
//go:generate godocdown -output=hue/lights/README.md hue/lights
//go:generate godocdown -output=hue/message/README.md hue/message

// Generate mocks.
//go:generate go get github.com/golang/mock/mockgen
//go:generate mockgen -source hue/hue.go -destination hue/mockHue/mockHue.go -package mockHue

// Disco Dance Party is program used for controlling Hue lights and run animations.
package main

func main() {}
