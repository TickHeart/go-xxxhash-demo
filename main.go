package main

import (
	"awesomeProject/cache"
	"awesomeProject/fileHash"
	"awesomeProject/utils"

	"github.com/kataras/iris/v12"
)

func main() {
	utils.LogSuc("create web server")
	app := iris.Default()
	router(app)
	err := app.Listen(":8080")
	if err != nil {
		return
	}
}

func router(app *iris.Application) {
	app.Get("/watcher", watcher)
	app.Get("/write", fileHash.Write)
}

func watcher(ctx iris.Context) {
	utils.LogSuc("create yaml hash server")
	go fileHash.WatcherInit(cache.MapCache)
	_, err := ctx.JSON(iris.Map{
		"name": "asdas",
	})

	if err != nil {
		return
	}
}
