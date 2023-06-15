package main

import (
	"final-project/alidada/app"
	"log"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start(":" + "3000")) // todo change port based on configs
}
