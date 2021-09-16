package internal

import (
	"context"
	"log"

	"github.com/ark-go/codefindbot/internal/jt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func startBotSession(ctx context.Context, client *telegram.Client) func(ctx context.Context) error {
	//	auth.NewFlow(auth.UserAuthenticator{})
	return func(ctx context.Context) error {
		log1, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		// Проверка статуса авторизации.
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return err
		}
		// Может быть уже аутентифицирован, если у нас есть действующий сеанс
		// в хранилище сессий.
		if !status.Authorized {
			// В противном случае выполните аутентификацию бота..
			log1.Info("Авторизация !!!")
			if _, err := client.Auth().Bot(ctx, jt.SecretEnv.BOT_TOKEN); err != nil {
				log1.Info("Авторизация ОШИБКА !!!" + err.Error())
				return err
			}
		}
		// Все в порядке, аутентификация вручную.
		log1.Info("Сеанс бота запущен")

		log.Println(".")
		//return telegram.RunUntilCanceled(ctx, client)
		// user := &tg.InputUserBox{
		// 	InputUser: &tg.InputUser{UserID: 22222222},
		// }
		//log.Println(">>>>", user)

		fulluser, err := client.API().UsersGetFullUser(ctx, &tg.InputUser{UserID: jt.SecretEnv.TgMoi1}) //
		if err == nil {
			log.Println(">>", fulluser.User.(*tg.User).AccessHash)
		} else {
			log.Println(">>>", err.Error())
		}

		<-ctx.Done() ///завершение
		//telegram.RunUntilCanceled(ctx, client)
		log1.Info("Сеанс бота закончен: " + ctx.Err().Error())
		return nil
	}
}
