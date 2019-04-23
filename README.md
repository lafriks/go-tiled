# go-tiled

[![GoDoc](https://godoc.org/github.com/lafriks/go-tiled?status.svg)](https://godoc.org/github.com/lafriks/go-tiled)
[![Build Status](https://cloud.drone.io/api/badges/lafriks/go-tiled/status.svg)](https://cloud.drone.io/lafriks/go-tiled)

Go library to parse Tiled map editor file format (TMX) and render map to image

Curently supports only orthogonal rendering

## Installing

    $ go get github.com/lafriks/go-tiled

You can use `go get -u` to update the package.

## Basic Usage:

```go
package main

import (
	"fmt"
	"os"

	"github.com/lafriks/go-tiled"
)

const mapPath = "maps/map.tmx" // path to your map

func main() {
    // parse tmx file
	gameMap, err := tiled.LoadFromFile(mapPath)

	if err != nil {
		fmt.Println("Error parsing map")
		os.Exit(2)
	}

	fmt.Print(gameMap)
}

```

## Documentation

For docs, see https://godoc.org/github.com/lafriks/go-tiled or run:

    $ godoc github.com/lafriks/go-tiled

