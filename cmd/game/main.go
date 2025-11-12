package main

import (
	"log"
	"path/filepath"

	"example.com/go-ogl-engine/Projects/daylight"
	"example.com/go-ogl-engine/internal/engine"
)

func main() {
	assets := filepath.Join("assets")
	e, err := engine.New(1280, 720, "Go OGL Engine", assets)
	if err != nil {
		log.Fatalln(err)
	}
	e.UseProject(daylight.New())
	e.Run()
}
