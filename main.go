package main

import (
	"awesomeProject/scheduler"
	"awesomeProject/utils"
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	watcherDirs("./hahaha", watcher)
	if err != nil {
		log.Fatal(err)
	}
	startWatch(watcher)
	defer watcher.Close()
}

func startWatch(watcher *fsnotify.Watcher) {
	var schedulerEvent *scheduler.Event
	for {
		select {
		case event, ok := <-watcher.Events:
			log.Printf("%v, %v", event, ok)
			if !ok {
				return
			}
			schedulerEvent = &scheduler.Event{
				Name: event.Op.String(),
				Path: event.Name,
			}
			signal := schedulerEvent.IgnoreSignal()

			if signal {
				schedulerEvent.ModelSchedulerModel(watcher)
			}

			log.Println(event.String())
		case err, ok := <-watcher.Errors:
			log.Printf("err: %v, %v", err, ok)
			if !ok {
				log.Println("error:", err)
				return
			}
		}
	}
}

func watcherDirs(dirPth string, watcher *fsnotify.Watcher) {
	dirs := utils.GetDirs(dirPth)
	for _, path := range dirs {
		err := watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}
}
