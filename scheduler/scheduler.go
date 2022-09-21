package scheduler

import (
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/fsnotify/fsnotify"
)

type Event struct {
	Path string
	Name string
	mu   sync.RWMutex
}

func (e *Event) IgnoreSignal() bool {
	models := []string{"CREATE", "REMOVE", "WRITE"}
	isModel := false
	isPath := false

	for _, val := range models {
		if val == e.Name {
			isModel = true
			break
		}
	}

	matchString, _ := regexp.MatchString("~", e.Path)
	if !matchString {
		isPath = true
	}

	return isModel && isPath
}

func (e *Event) ModelSchedulerModel(watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	switch e.Name {
	case "CREATE":
		go e.createModel(watcher, wg)
	case "REMOVE":

	case "WRITE":
		go e.writeModel(watcher, wg)
	}
}

func (e *Event) createModel(watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	e.mu.Lock()
	defer e.mu.Unlock()
	defer wg.Done()
	ms, _ := regexp.MatchString("yaml", e.Path)
	isHash, _ := regexp.MatchString("\\.[\\d]+\\.", e.Path)
	if isHash {
		return
	}
	if ms {
		file, _ := os.ReadFile(e.Path)
		body := string(file)
		//if len(body) <= 0 {
		//	return
		//}
		fileRenameToXXHash(body, e.Path)
	} else {
		err := watcher.Add(e.Path)
		if err != nil {
			return
		}
	}
}
func (e *Event) writeModel(watcher *fsnotify.Watcher, wg *sync.WaitGroup) {
	e.mu.Lock()
	defer e.mu.Unlock()
	defer wg.Done()
	ms, _ := regexp.MatchString("yaml", e.Path)
	if ms {
		file, _ := os.ReadFile(e.Path)
		body := string(file)
		//if len(body) <= 0 {
		//	return
		//}
		updateFilenameXXHash(body, e.Path)
	}
}

func fileRenameToXXHash(body string, path string) {
	sum64String := xxhash.Sum64String(body)
	compile := regexp.MustCompile("yaml")
	allString := compile.ReplaceAllString(path, "")
	allString = allString + strconv.FormatUint(sum64String, 10) + ".yaml"
	err := os.Rename(path, allString)
	if err != nil {
		return
	}
}

func updateFilenameXXHash(body string, path string) {
	sum64String := xxhash.Sum64String(body)
	compile := regexp.MustCompile(`\\.[\\d]+\\.`)
	allString := compile.ReplaceAllString(path, "."+strconv.FormatUint(sum64String, 10)+".")
	err := os.Rename(path, allString)
	if err != nil {
		return
	}
}
