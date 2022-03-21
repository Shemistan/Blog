package main

import (
	"context"
	"flag"
	"log"

	"github.com/Shemistan/Blog/internal/app"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "path", "configs/configs.toml", "path to config file in .toml format")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	a, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatal("failed to create app")
	}

	err = a.Run()
	if err != nil {
		log.Fatal("failed to run app")
	}
}
