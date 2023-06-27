package main

import (
	"log"
	"qsms/app"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start(":" + "3000"))
}
