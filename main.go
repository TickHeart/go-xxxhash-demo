package main

import (
	"awesomeProject/fileHash"
	"awesomeProject/utils"
	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
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
	app.Get("/write", write)
}

func watcher(ctx iris.Context) {
	utils.LogSuc("create yaml hash server")
	go fileHash.WatcherInit()
	_, err := ctx.JSON(iris.Map{
		"name": "asdas",
	})

	if err != nil {
		return
	}
}

type Y struct {
	Mode string `yaml:"mode"`
	List []User `yaml:"list"`
}
type User struct {
	Name string `yaml:"name"`
	Age  string `yaml:"age"`
}

type WriteEvent struct {
	wg        sync.WaitGroup
	rw        sync.RWMutex
	y         interface{}
	username  string
	tableName string
}

func (w *WriteEvent) setInfo(username, tableName string, y interface{}) {
	w.username = username
	w.tableName = tableName
	w.y = y
}

func write(ctx iris.Context) {
	username := ctx.URLParam("username")
	tableName := ctx.URLParam("tableName")
	y := Y{
		Mode: "update",
		List: []User{{
			"wuhongbin",
			"15",
		}},
	}
	var we WriteEvent
	we.setInfo(username, tableName, y)
	files := we.CheckHashTableFile()
	we.updateUserTableYaml(files)
}

func (w *WriteEvent) updateUserTableYaml(files []string) {
	marshal, _ := yaml.Marshal(w.y)
	if len(files) >= 1 {
		filePath := files[0]
		err := os.WriteFile(filePath, marshal, 0666)
		if err != nil {
			return
		}
		fileHash.UpdateFilenameXXHash(string(marshal), filePath, false)
	} else {
		_ = ioutil.WriteFile("/Users/wuhongbin/Desktop/awesomeProject/hahaha/hash/"+w.tableName+".yaml", marshal, 0666)
		fileHash.FileRenameToXXHash(string(marshal), "/Users/wuhongbin/Desktop/awesomeProject/hahaha/hash/"+w.tableName+".yaml")
	}
}

func (w *WriteEvent) CheckHashTableFile() []string {
	files := utils.GetYamlFiles("/Users/wuhongbin/Desktop/awesomeProject/hahaha/hash")
	fileChan := make(chan string, len(files))
	for i := 0; i < len(files); i++ {
		w.wg.Add(1)
		val := files[i]
		go func() {
			defer w.wg.Done()
			baseFileName := w.username + "\\." + w.tableName
			matchString, _ := regexp.MatchString(baseFileName+"\\.\\d+?\\.yaml", val)
			log.Println(baseFileName + "\\.\\d+?\\.yaml")
			if matchString {
				fileChan <- val
			}
		}()
	}
	w.wg.Wait()
	var fis []string
	close(fileChan)
	for file := range fileChan {
		strings := append(fis, file)
		fis = strings
	}
	return files
}
