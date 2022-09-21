package fileHash

import (
	"github.com/cespare/xxhash/v2"
	"github.com/fsnotify/fsnotify"
	"os"
	"regexp"
	"strconv"
)

func ModelSchedulerModel(eventPath, eventName string, watcher *fsnotify.Watcher) {
	switch eventName {
	case "CREATE":
		createModel(eventPath, eventName, watcher)
		break
	case "REMOVE":

		break
	case "WRITE":
		writeModel(eventPath, eventName, watcher)
		break
	}
}

func createModel(eventPath, eventName string, watcher *fsnotify.Watcher) {
	ms, _ := regexp.MatchString("yaml", eventPath)
	isHash, _ := regexp.MatchString("\\.[\\d]+\\.", eventPath)
	if isHash {
		return
	}
	if ms {
		file, _ := os.ReadFile(eventPath)
		body := string(file)
		//if len(body) <= 0 {
		//	return
		//}
		FileRenameToXXHash(body, eventPath)
	} else {
		err := watcher.Add(eventPath)
		if err != nil {
			return
		}
	}
}
func writeModel(eventPath, eventName string, watcher *fsnotify.Watcher) {
	ms, _ := regexp.MatchString("yaml", eventPath)
	if ms {
		file, _ := os.ReadFile(eventPath)
		body := string(file)
		//if len(body) <= 0 {
		//	return
		//}
		UpdateFilenameXXHash(body, eventPath, false)
	}
}

func FileRenameToXXHash(body string, path string) {
	sum64String := xxhash.Sum64String(body)
	compile := regexp.MustCompile("yaml")
	allString := compile.ReplaceAllString(path, "")
	allString = allString + strconv.FormatUint(sum64String, 10) + ".yaml"
	err := os.Rename(path, allString)
	if err != nil {
		return
	}
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
