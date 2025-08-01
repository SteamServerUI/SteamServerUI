package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func InitPlugins() {
	pluginDir := "./UIMod/plugins"

	files, err := os.ReadDir(pluginDir)
	if err != nil {
		panic(err)
	}

	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" {
			fullPath := filepath.Join(pluginDir, file.Name())
			fmt.Printf("Lade Plugin: %s\n", file.Name())

			_, err := i.EvalPath(fullPath)
			if err != nil {
				fmt.Printf("Fehler in %s: %v\n", file.Name(), err)
				continue
			}

			v, err := i.Eval("main.Run()")
			if err != nil {
				fmt.Printf("Fehler beim Ausführen von Run(): %v\n", err)
			}
			fmt.Println("Returned:", v)
		}
	}
}
