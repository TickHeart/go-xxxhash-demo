package fileHash

import (
	"awesomeProject/utils"
	"github.com/cespare/xxhash/v2"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var w sync.WaitGroup

type Event struct {
	EventPath string
	EventName string
	wrs       sync.RWMutex
}

func (receiver *Event) SetInfo(event fsnotify.Event) {
	receiver.EventName = event.Op.String()
	receiver.EventName = event.Name
}

func (receiver *Event) IgnoreSignal() bool {
	models := []string{"CREATE", "REMOVE", "WRITE"}
	isModel := false
	isPath := false

	for i := 0; i < len(models); i++ {
		val := models[i]
		if val == receiver.EventName {
			isModel = true
		}
	}
	matchString, _ := regexp.MatchString("~", receiver.EventPath)
	if !matchString {
		isPath = true
	}

	return isModel && isPath
}

func (receiver *Event) inRLock(fn func() interface{}) interface{} {
	receiver.wrs.RLock()
	i := fn()
	receiver.wrs.RUnlock()
	return i
}
func (receiver *Event) inLock(fn func() interface{}) interface{} {
	receiver.wrs.Lock()
	i := fn()
	receiver.wrs.Unlock()
	return i
}

func WatcherInit() {
	initialize("/Users/wuhongbin/Desktop/awesomeProject/hahaha")
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
				e := Event{}
				e.SetInfo(event)
				signal := e.IgnoreSignal()

				if signal {
					ModelSchedulerModel(&e, watcher)
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

type YamlHash struct {
	hash     string
	fileName string
}

func initialize(dirPth string) {
	utils.LogInfo("初始化开始")
	yamlFiles := utils.GetYamlFiles(dirPth)

	errYamlFiles := make(chan YamlHash, len(yamlFiles))

	checkYamlHash(yamlFiles, errYamlFiles)

	amendmentYamlFilesHash(errYamlFiles)

	utils.LogInfo("检查结束")
}

func checkYamlHash(yamlFiles []string, errYamlFiles chan YamlHash) {

	if len(yamlFiles) <= 0 {
		return
	}

	for i := 0; i < len(yamlFiles); i++ {
		w.Add(1)
		file := yamlFiles[i]
		go func() {
			defer w.Done()
			fileBody, _ := os.ReadFile(file)
			sum64String := xxhash.Sum64String(string(fileBody))
			matchString, _ := regexp.MatchString(strconv.FormatUint(sum64String, 10), file)

			if !matchString {
				errYamlFiles <- YamlHash{
					hash:     strconv.FormatUint(sum64String, 10),
					fileName: file,
				}
			}
		}()

	}
}

func amendmentYamlFilesHash(errYamlFiles chan YamlHash) {
	w.Wait()
	close(errYamlFiles)
	id := 1
	for iv := range errYamlFiles {
		w.Add(1)
		utils.LogInfo("哈西出现错误 " + iv.fileName)
		go func(yamlInfo YamlHash) {
			UpdateFilenameXXHash(yamlInfo.hash, yamlInfo.fileName, true)
			utils.LogInfo("修改完成个数" + strconv.Itoa(id) + "：" + yamlInfo.fileName)
			id += 1
			w.Done()
		}(iv)
	}
	w.Wait()
	utils.LogInfo("所有文件 hash 修改成功")
}
