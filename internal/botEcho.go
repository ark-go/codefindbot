package internal

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

var sender *message.Sender

func StartBotEcho() {
	Run(func(ctx context.Context, log *zap.Logger) error {
		// Dispatcher handles incoming updates.
		dispatcher := tg.NewUpdateDispatcher()
		opts := telegram.Options{
			//Logger:        log,
			UpdateHandler: dispatcher,
		}

		return telegram.BotFromEnvironment(ctx, opts, func(ctx context.Context, client *telegram.Client) error {
			// Raw MTProto API client, allows making raw RPC calls.
			api := tg.NewClient(client)

			// Помощник для отправки сообщений.
			sender = message.NewSender(api)

			// Setting up handler for incoming message.
			dispatcher.OnNewMessage(onNewMessage)

			return nil
		}, telegram.RunUntilCanceled)
	})
}
