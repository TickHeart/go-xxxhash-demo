package fileHash

import (
	"awesomeProject/utils"
	"github.com/fsnotify/fsnotify"
	"log"
	"regexp"
)

func WatcherInit() {
	watcher, err := fsnotify.NewWatcher()
	watcherDirs("/Users/wuhongbin/Desktop/awesomeProject/hahaha", watcher)
	if err != nil {
		log.Fatal(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		log.Println("监听结束")
		err := watcher.Close()
		if err != nil {

		}
	}(watcher)
	err = watcher.Add("/Users/wuhongbin/Desktop/awesomeProject/hahaha")
	log.Println("确认已经开始监听")
	func() {
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
					ModelSchedulerModel(eventPath, eventName, watcher)
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

	if err != nil {
		log.Fatal(err)
	}
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
