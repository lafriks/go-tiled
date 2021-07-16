package main

import (
	"fmt"
	"os"

	"github.com/lafriks/go-tiled"
)

const (
	tmxPath = "../go-tiled/assets/test_wangsets_w_properties_map.tmx"
)

func main() {
	var tiledMap *tiled.Map
	var err error
	tiledMap, err = tiled.LoadFromFile(tmxPath)
	if err != nil {
		fmt.Printf("error parsing gameMap: %s", err.Error())
		os.Exit(2)
	}

	// get the tile so that we can read its data
	tile, err := tiledMap.Layers[0].Tiles[0].Tileset.GetTilesetTile(15)

	fmt.Print(tile)

}
