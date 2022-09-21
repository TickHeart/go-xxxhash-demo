package main

import (
	"awesomeProject/fileHash"
	"awesomeProject/utils"
	"github.com/kataras/iris/v12"
)

func main() {
	utils.LogSuc("create web server")
	app := iris.Default()
	app.Get("/someGet", getting)
	err := app.Listen(":8080")
	if err != nil {
		return
	}
}

func getting(ctx iris.Context) {
	utils.LogSuc("create yaml hash server")
	go fileHash.WatcherInit()
	_, err := ctx.JSON(iris.Map{
		"name": "asdas",
	})

	if err != nil {
		return
	}
}
