package jt

import (
	"log"
	"os"
	"path/filepath"
)

var rootDir string
var err error

func init() {
	rootDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln("не определить рабочий каталог")
		rootDir = "."
	}
	_ = rootDir
}
