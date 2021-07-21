package main

import (
	"fmt"
	"os"

	"github.com/lafriks/go-tiled"
)

const (
	wangTmxPath = "../go-tiled/assets/test_wangsets_map.tmx"
)

func main() {
	var tiledMap *tiled.Map
	var err error
	tiledMap, err = tiled.LoadFromFile(wangTmxPath)
	if err != nil {
		fmt.Printf("error parsing tiledMap: %s", err.Error())
		os.Exit(2)
	}

	fmt.Print(tiledMap.Tilesets[0].WangSets)

	wangColors, err := tiledMap.Tilesets[0].WangSets[0].GetWangColors(16)
	fmt.Print(wangColors)
}
