package main

import (
	"github.com/Veerse/podcast-feed-api/app"
	"github.com/Veerse/podcast-feed-api/config"
)

func main() {
	a := app.App{}
	a.Initialize(config.GetConfigFromFile())
	a.Run()
}