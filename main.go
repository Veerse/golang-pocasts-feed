package main

import (
	"github.com/Veerse/podcast-feed-api/app"
	"github.com/Veerse/podcast-feed-api/config"
	"log"
)

func main() {
	a := app.App{}

	if c, err := config.GetConfigFromFile(); err != nil {
		panic(err)
	} else {
		if err := a.Initialize(c); err != nil {
			log.Fatalf("initialization : %s", err.Error())
		}
	}

	a.Run()
}
