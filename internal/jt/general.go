package jt

import (
	logorig "log"
	"os"
	"path/filepath"
)

var log *logorig.Logger
var RootDir string
var err error
var Infocod string

type secretEnv struct {
	IsSecret  string
	App_id    int
	App_hash  string
	BOT_TOKEN string
	TgMoi1    int
	TgMoi2    int
	TelMoi    string
}

var SecretEnv *secretEnv

func init() {
	RootDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("не определить рабочий каталог")
		RootDir = "."
	}
	_ = RootDir
	SecretEnv = &secretEnv{}
}

func init() {
	log = logorig.New(os.Stdout, "- ", logorig.LstdFlags)
}

// log, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel), zap.AddStacktrace(zapcore.FatalLevel))
//"go.uber.org/zap"
//"go.uber.org/zap/zapcore"
