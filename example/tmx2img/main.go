package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

func main() {
	flag.Parse()

	filename := flag.Arg(0)
	img := flag.Arg(1)
	if img == "" {
		img = "map.png"
	}

	m, err := tiled.LoadFromFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	rend, err := render.NewRenderer(m)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = rend.RenderVisibleLayers(); err != nil {
		fmt.Println(err)
		return
	}
	//rend.RenderLayer(1)

	w, err := os.Create(img)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer w.Close()
	rend.SaveAsPng(w)
}
