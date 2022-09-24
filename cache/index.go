package cache

import (
	"os"
	"path"
	"regexp"
)

type Cache struct {
	fileName string
	fileBody string
	fileHash string
	fileFlag string
}

type MapCacheType = map[string]Cache

var MapCache = make(MapCacheType)

func SetMapCache(yamlPath string) {
	fileName := path.Base(yamlPath)
	fileBat, _ := os.ReadFile(yamlPath)
	fileBody := string(fileBat)
	infos := hashFileNameToSplitInfo(fileName)
	if len(infos) <= 0 {
		return
	}
	fileFlag := infos[0] + "&" + infos[1]

	cache := Cache{
		fileName: fileName,
		fileBody: fileBody,
		fileHash: infos[2],
		fileFlag: fileFlag,
	}
	MapCache[fileFlag] = cache
}

func hashFileNameToSplitInfo(fileName string) []string {
	compile := regexp.MustCompile("(.+?)\\.(.+?)\\.(\\d+?)\\.yaml")
	matchString := compile.FindStringSubmatch(fileName)

	if len(matchString) <= 0 {
		return []string{}
	}
	return []string{matchString[1], matchString[2], matchString[3]}
}
