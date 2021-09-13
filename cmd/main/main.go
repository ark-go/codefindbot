package main

import (
	"log"

	"github.com/ark-go/codefindbot/internal/jt"
)

var versionProg string

type secretEnv struct {
	IsSecret  string
	APP_ID    string
	APP_HASH  string
	BOT_TOKEN string
}

var SecretEnv *secretEnv

func init() {
	SecretEnv = &secretEnv{}
}

func main() {
	log.Println("Version:", versionProg)
	jt.LoadEnvironment(SecretEnv)
	if SecretEnv.IsSecret != "" {
		log.Println("Загружены данные из secret.env")
		if SecretEnv.IsSecret == "1" {
			jt.PrintEnvironment(SecretEnv)
		}
	}

}
