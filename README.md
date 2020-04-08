# go-tiled

[![GoDoc](https://godoc.org/github.com/lafriks/go-tiled?status.svg)](https://godoc.org/github.com/lafriks/go-tiled)
[![Build Status](https://cloud.drone.io/api/badges/lafriks/go-tiled/status.svg)](https://cloud.drone.io/lafriks/go-tiled)

Go library to parse Tiled map editor file format (TMX) and render map to image. Currently supports only orthogonal rendering out-of-the-box.

## Installing

    $ go get github.com/lafriks/go-tiled

You can use `go get -u` to update the package. You can also just import and start using the package directly if you're using Go modules, and Go will then download the package on first compilation.

## Basic Usage:

```go
package main

import (
	"fmt"
	"os"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

const mapPath = "maps/map.tmx" // Path to your Tiled Map.

func main() {
    // Parse .tmx file.
	gameMap, err := tiled.LoadFromFile(mapPath)

	if err != nil {
		fmt.Println("Error parsing map")
		os.Exit(2)
	}

	fmt.Println(gameMap)

	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer := render.NewRenderer(gameMap)

	// Render just layer 0 to the Renderer.
	renderer.RenderLayer(0)

	// Get a reference to the Renderer's output, an image.NRGBA struct.
	img := renderer.Result

	// Clear the render result after copying the output if separation of 
	// layers is desired.
	renderer.Clear()

	// And so on. You can also export the image to a file by using the
	// Renderer's Save functions.

}

```

## Documentation

For further documentation, see https://pkg.go.dev/github.com/lafriks/go-tiled or run:

    $ godoc github.com/lafriks/go-tiled

