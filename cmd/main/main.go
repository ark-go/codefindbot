package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/ark-go/codefindbot/internal"

	"github.com/ark-go/codefindbot/internal/jt"
)

var versionProg string

func init() {
	//	jt.SecretEnv = &secretEnv{}
}

func main() {

	log.Println("Version:", versionProg)
	jt.LoadEnvironment(jt.SecretEnv)
	if jt.SecretEnv.IsSecret != "" {
		log.Println("Загружены данные из secret.env")
		if jt.SecretEnv.IsSecret == "1" {
			jt.PrintEnvironment(jt.SecretEnv)
		}
	}
	os.Setenv("BOT_TOKEN", jt.SecretEnv.BOT_TOKEN)
	os.Setenv("APP_ID", strconv.Itoa(jt.SecretEnv.App_id))
	os.Setenv("APP_HASH", jt.SecretEnv.App_hash)
	//internal.StartBotEcho()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		//internal.RunBot()
		internal.StartBot()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		//internal.StartClientTg()
	}()
	//time.Sleep(time.Millisecond * 1)
	wg.Wait()
	log.Println("Завершение работы")

}

/*
defer func() {
		signal.Stop(c)
		cancel()
	}()
*/
