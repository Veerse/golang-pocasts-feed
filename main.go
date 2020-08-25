package main

import (
	"github.com/Veerse/podcast-feed-api/app"
	"github.com/Veerse/podcast-feed-api/config"
)

func main() {
	a := app.App{}

	if err := a.Initialize(config.GetConfigFromFile()); err != nil {
		panic(err)
	}

	a.Run()
}