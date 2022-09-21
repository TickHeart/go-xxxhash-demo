package main

import (
	"awesomeProject/scheduler"
	"awesomeProject/utils"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	watcherDirs("./hahaha", watcher)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1000)
	startWatch(watcher, &wg)
	wg.Wait()
}

func startWatch(watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	var schedulerEvent *scheduler.Event
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				log.Printf("%v", event)
				if !ok {
					log.Fatalln("接收失败")
					return
				}
				schedulerEvent = &scheduler.Event{
					Name: event.Op.String(),
					Path: event.Name,
				}
				signal := schedulerEvent.IgnoreSignal()

				if signal {
					schedulerEvent.ModelSchedulerModel(watcher, wg)
				}

				log.Println(event.String())
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("error:", err)
					return
				}
			}
		}
	}()
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			panic(err)
		}
	}(watcher)
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
