package fileHash

import (
	"awesomeProject/constants"
	"awesomeProject/utils"
	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
)

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
	rr        sync.Mutex
	y         interface{}
	username  string
	tableName string
}

func (w *WriteEvent) setInfo(username, tableName string, y interface{}) {
	w.username = username
	w.tableName = tableName
	w.y = y
}
func (w *WriteEvent) updateUserTableYaml(files []string) {
	marshal, _ := yaml.Marshal(w.y)
	standard := ""
	if len(files) >= 1 {
		standard = files[0]
		err := os.WriteFile(standard, marshal, 0666)
		if err != nil {
			return
		}
	} else {
		standard = constants.GetHashPath() + "/" + w.username + "." + w.tableName + ".yaml"
		_ = ioutil.WriteFile(standard, marshal, 0666)
	}
	//FileRenameToXXHash(string(marshal), standard)
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
			matchString, _ := regexp.MatchString(baseFileName+"(\\.\\d+?)?\\.yaml", val)
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

	return fis
}

var f = make(map[string]func())

func Write(ctx iris.Context) {

	username := ctx.URLParam("username")
	tableName := ctx.URLParam("tableName")
	y := Y{
		Mode: "update",
		List: []User{{
			"123123",
			"15",
		}},
	}
	var we WriteEvent
	we.setInfo(username, tableName, y)

	if f[username+"&"+tableName] == nil {
		f[username+"&"+tableName] = func() {
			var l sync.Mutex
			l.Lock()
			files := we.CheckHashTableFile()
			we.updateUserTableYaml(files)
			l.Unlock()
		}

	}
	f[username+"&"+tableName]()

	_, err := ctx.JSON(y)
	if err != nil {
		return
	}
}
