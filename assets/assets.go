package assets

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var gifs = sync.Map{}

type Asset struct {
	Path string
	File os.FileInfo
}

func Load() {
	items := List()

	for _, item := range items {

		name := item.File.Name()

		if item.File.IsDir() || strings.TrimSpace(name) == "" || !strings.Contains(name, ".gif") {
			continue
		}

		content, err := ioutil.ReadFile(item.Path)
		if err != nil {
			log.Panicln(err)
		}

		key := RemoveExt(item.File.Name())
		gifs.Store(key, content)
	}
}

func RemoveExt(fileName string) string {
	index := strings.Index(fileName, ".")
	return fileName[:index]
}

func List() []*Asset {
	var assetList []*Asset

	root := "."

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		assetList = append(assetList, &Asset{Path: path, File: info})
		return err
	})
	if err != nil {
		panic(err)
	}

	return assetList
}

func GetNames() []string {
	var names []string
	gifs.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

func Count() int {
	var c int
	gifs.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func Get(name string) []byte {
	b, ok := gifs.Load(name)
	if ok {
		return b.([]byte)
	}
	return []byte{}
}
