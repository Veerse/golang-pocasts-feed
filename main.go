package main

import (
	"github.com/Veerse/podcast-feed-api/app"
	"github.com/Veerse/podcast-feed-api/config"
)

func main() {
	a := app.App{}

	if config, err := config.GetConfigFromFile(); err != nil {
		panic(err)
	} else {
		if err := a.Initialize(config); err != nil {
			panic(err)
		}
	}

	a.Run()
}
