package main

import (
	"github.com/Veerse/podcast-feed-api/app"
	"github.com/Veerse/podcast-feed-api/config"
)

func main() {
	a := app.App{}

	if c, err := config.GetConfigFromFile(); err != nil {
		panic(err)
	} else {
		if err := a.Initialize(c); err != nil {
			panic(err)
		}
	}

	a.Run()
}
