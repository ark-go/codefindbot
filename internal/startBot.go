package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ark-go/codefindbot/internal/jt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// Grab those from https://my.telegram.org/apps.
// appID := flag.Int("api-id", 0, "app id")
// appHash := flag.String("api-hash", "hash", "app hash")
// // Get it from bot father.
// token := flag.String("token", "", "bot token")
// flag.Parse()

func StartBot() {
	//	log, err := zap.NewDevelopment()
	log1, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log1.Sync() }()
	//ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		//  os.Exit(1)
	}()
	//	go func() {
	if err := startClientBot(ctx); err != nil {
		log1.Fatal("Run failed", zap.Error(err))
	}
	log1.Info("Выход")
	//	}()

	// Done.
}
func startClientBot(ctx context.Context) error {
	log1, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	dispatcher := tg.NewUpdateDispatcher()
	//log.Info("::::::::::::::::::::::::", jt.SecretEnv.App_id, jt.SecretEnv.App_hash)
	client := telegram.NewClient(jt.SecretEnv.App_id, jt.SecretEnv.App_hash, telegram.Options{
		SessionStorage: sessionStorage.New("sessionBot.dat"),
		Logger:         log1,
		UpdateHandler:  dispatcher,
		// Middlewares: []telegram.Middleware{
		// 	utilites.PrettyMiddleware(),
		// },
	})
	// Создаем клиента
	api := tg.NewClient(client)
	// Помощник для отправки сообщений.
	sender = message.NewSender(api)
	// Настройка обработчика входящего сообщения.
	dispatcher.OnNewMessage(onNewMessage)

	errc := client.Run(ctx, startBotSession(ctx, client)) // Запускает сеанс бота

	return errc
}
