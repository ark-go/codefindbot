package internal

import (
	"context"
	"log"

	"github.com/ark-go/codefindbot/internal/jt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func StartClientTg() {
	dispatcher := tg.NewUpdateDispatcher()
	client := telegram.NewClient(jt.SecretEnv.App_id, jt.SecretEnv.App_hash, telegram.Options{
		SessionStorage: sessionStorage.New("sessionClient.dat"),
		//Logger:         log,
		UpdateHandler: dispatcher,
	})
	phone := jt.SecretEnv.TelMoi
	flow := auth.NewFlow(
		termAuth{phone: phone},
		auth.SendCodeOptions{},
	)
	// Создаем клиента
	api := tg.NewClient(client)
	// Помощник для отправки сообщений.
	sender = message.NewSender(api)
	// Настройка обработчика входящего сообщения.
	dispatcher.OnNewMessage(onNewMessage)

	if err := client.Run(context.Background(), func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		log.Println("Success")
		// It is only valid to use client while this function is not returned
		// and ctx is not cancelled.
		// api := client.API()
		// _ = api
		// Now you can invoke MTProto RPC requests by calling the API.
		// ...

		// client.SendMessage(ctx, &tg.MessagesSendMessageRequest{
		// 	Peer:    &tg.InputPeerChat{ChatID: 222222}, // tg.InputPeerClass & tg.InputUser{UserID: 222222}, //.() //    User{UserID: 222222},
		// 	Message: "Хай",
		// })

		message.NewSender(tg.NewClient(client)).Self().Text(ctx, "Hi777!") // в избранное отправил
		//	message.Sender()
		// _, err := sender.Resolve("222222").Text(ctx, "ХААЙ!")
		// if err != nil {
		// 	log.Println("send", err.Error())
		// }

		_, err := sender.To(&tg.InputPeerUser{UserID: jt.SecretEnv.TgMoi2}).Text(ctx, "Привет3333")
		if err != nil {
			log.Println("send", err.Error())
		}
		log.Println("Зaгрузились")
		// Return to close client connection and free up resources.
		<-ctx.Done()
		return nil
	}); err != nil {
		panic(err)
	}
	// Client is closed.

}

/*
 g, err := tg.DecodePeer(buf)
 if err != nil {
     panic(err)
 }
 switch v := g.(type) {
 case *tg.PeerUser: // peerUser#9db1bc6d
 case *tg.PeerChat: // peerChat#bad0e5bb
 case *tg.PeerChannel: // peerChannel#bddde532
 default: panic(v)
 }
*/
