package qsms

import (
	"log"
	"qsms/app"
)

func main() {
	app := app.NewApp()
	log.Fatalln(app.Start(":" + "3000")) // todo change port based on configs
}
