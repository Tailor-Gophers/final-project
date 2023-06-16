package main

import (
	"alidada/app"
	"log"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start(":" + "3000")) // todo change port based on configs
}
