package fileHash

import (
	"github.com/cespare/xxhash/v2"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func ModelSchedulerModel(e *Event, watcher *fsnotify.Watcher) {
	switch e.EventName {
	case "CREATE":
		createModel(e, watcher)
		break
	case "REMOVE":

		break
	case "WRITE":
		writeModel(e)
		break
	}
}

func createModel(e *Event, watcher *fsnotify.Watcher) {
	ms, _ := regexp.MatchString("yaml", e.EventPath)
	isHash, _ := regexp.MatchString("\\.[\\d]+\\.", e.EventPath)
	if isHash {
		return
	}
	if ms {
		lock := e.inRLock(func() interface{} {
			file, _ := os.ReadFile(e.EventPath)
			body := string(file)
			return body
		})
		e.inLock(func() interface{} {
			FileRenameToXXHash(lock.(string), e.EventPath)
			return true
		})
	} else {
		err := watcher.Add(e.EventPath)
		if err != nil {
			return
		}
	}
}

func writeModel(e *Event) {
	ms, _ := regexp.MatchString("yaml", e.EventPath)
	if ms {
		lock := e.inRLock(func() interface{} {
			file, _ := os.ReadFile(e.EventPath)
			body := string(file)
			return body
		})
		e.inLock(func() interface{} {
			UpdateFilenameXXHash(lock.(string), e.EventPath, false)
			return true
		})

	}
}

func FileRenameToXXHash(body string, path string) string {
	sum64String := xxhash.Sum64String(body)
	compile := regexp.MustCompile("yaml")
	allString := compile.ReplaceAllString(path, "")
	allString = allString + strconv.FormatUint(sum64String, 10) + ".yaml"
	err := os.Rename(path, allString)
	if err != nil {
		return ""
	}
	_, file := filepath.Split(path)
	return filepath.Base(file)
}

func UpdateFilenameXXHash(body string, path string, isHash bool) {
	hash := ""
	if !isHash {
		sum64String := xxhash.Sum64String(body)
		hash = strconv.FormatUint(sum64String, 10)
	} else {
		hash = body
	}

	compile := regexp.MustCompile("\\.[\\d]+\\.")
	allString := compile.ReplaceAllString(path, "."+hash+".")
	err := os.Rename(path, allString)
	if err != nil {
		return
	}
}
