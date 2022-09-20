package main

import (
	"awesomeProject/scheduler"
	"awesomeProject/utils"
	"github.com/fsnotify/fsnotify"
	"log"
	"regexp"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	watcherDirs("./hahaha", watcher)
	if err != nil {
		log.Fatal(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {

		}
	}(watcher)
	cl := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				eventPath := event.Name
				eventName := event.Op.String()
				signal := ignoreSignal(eventPath, eventName)

				if signal {
					scheduler.ModelSchedulerModel(eventPath, eventName, watcher)
				}

				log.Println(event.String())
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./hahaha")
	if err != nil {
		log.Fatal(err)
	}
	<-cl
}

func ignoreSignal(eventPath, eventName string) bool {
	models := []string{"CREATE", "REMOVE", "WRITE"}
	isModel := false
	isPath := false

	for i := 0; i < len(models); i++ {
		val := models[i]
		if val == eventName {
			isModel = true
		}
	}
	matchString, _ := regexp.MatchString("~", eventPath)
	if !matchString {
		isPath = true
	}

	return isModel && isPath
}

func watcherDirs(dirPth string, watcher *fsnotify.Watcher) {
	dirs := utils.GetDirs(dirPth)
	for i := 1; i < len(dirs); i++ {
		val := dirs[i]
		err := watcher.Add(val)
		if err != nil {
			return
		}
	}
}
