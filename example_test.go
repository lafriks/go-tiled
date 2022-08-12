// Examples
package tiled_test

import (
	"fmt"
	"os"

	"github.com/lafriks/go-tiled"
)

func ExampleTileset_GetTilesetTile() {
	var tiledMap *tiled.Map
	var err error
	tiledMap, err = tiled.LoadFile("assets/test_wangsets_w_properties_map.tmx")
	if err != nil {
		fmt.Printf("error parsing tiledMap: %s", err.Error())
		os.Exit(2)
	}

	// get the tile so that we can read its data
	tile, err := tiledMap.Layers[0].Tiles[0].Tileset.GetTilesetTile(15)
	if err != nil {
		panic(err)
	}

	fmt.Print(tile.ID)

	// Output:
	// 15
}

func ExampleWangSet_GetWangColors() {
	var tiledMap *tiled.Map
	var err error
	tiledMap, err = tiled.LoadFile("assets/test_wangsets_map.tmx")
	if err != nil {
		fmt.Printf("error parsing tiledMap: %s", err.Error())
		os.Exit(2)
	}

	fmt.Println(tiledMap.Tilesets[0].WangSets[0].Name)

	wangColors, err := tiledMap.Tilesets[0].WangSets[0].GetWangColors(16)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  TopRight:    %s\n", wangColors[tiled.TopRight].Name)
	fmt.Printf("  BottomRight: %s\n", wangColors[tiled.BottomRight].Name)
	fmt.Printf("  BottomLeft:  %s\n", wangColors[tiled.BottomLeft].Name)
	fmt.Printf("  TopLeft:     %s\n", wangColors[tiled.TopLeft].Name)

	// Output:
	// Summer
	//   TopRight:    Rock
	//   BottomRight: Water
	//   BottomLeft:  Water
	//   TopLeft:     Rock
}
