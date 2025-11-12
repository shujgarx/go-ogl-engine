package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"example.com/go-ogl-engine/Projects/daylight"
	"example.com/go-ogl-engine/Projects/fallline"
	"example.com/go-ogl-engine/internal/engine"
)

var defaultProject = "fallline"

var projectFactories = map[string]func() engine.Project{
	"daylight": func() engine.Project { return daylight.New() },
	"fallline": func() engine.Project { return fallline.New() },
}

func projectNames() []string {
	names := make([]string, 0, len(projectFactories))
	for k := range projectFactories {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func main() {
	projectFlag := flag.String("project", "", "project name to run")
	listFlag := flag.Bool("list-projects", false, "list available projects and exit")
	flag.Parse()

	if *listFlag {
		fmt.Println("Available projects:")
		for _, name := range projectNames() {
			fmt.Printf(" - %s\n", name)
		}
		return
	}

	choice := strings.ToLower(strings.TrimSpace(*projectFlag))
	if choice == "" {
		if env := strings.TrimSpace(os.Getenv("ENGINE_PROJECT")); env != "" {
			choice = strings.ToLower(env)
		}
	}
	if choice == "" {
		choice = defaultProject
	}

	factory, ok := projectFactories[choice]
	if !ok {
		fmt.Printf("Unknown project \"%s\". Available projects: %s\n", choice, strings.Join(projectNames(), ", "))
		os.Exit(1)
	}

	assets := filepath.Join("assets")
	e, err := engine.New(1280, 720, "Go OGL Engine", assets)
	if err != nil {
		log.Fatalln(err)
	}
	e.UseProject(factory())
	fmt.Printf("Running project: %s\n", choice)
	e.Run()
}
