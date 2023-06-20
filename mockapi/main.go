package main

import (
	"log"
	"mockapi/app"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start(":" + "3001")) // todo change port based on configs
}
