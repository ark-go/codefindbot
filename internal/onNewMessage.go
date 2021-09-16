package internal

import (
	"context"
	"log"

	"github.com/gotd/td/tg"
)

func onNewMessage2(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
	//! https://core.telegram.org/constructor/updateNewMessage
	m, ok := u.Message.(*tg.Message) // *tg.MessageService
	if !ok || m.Out {
		//Исходящее сообщение, не интересно.
		return nil
	}
	log.Println("------------0000000000000000000000-------------")
	if peerUser, ok := m.FromID.(*tg.PeerUser); ok {
		log.Println("Послал id: ", peerUser.UserID)
	} else if _, ok := m.FromID.(*tg.PeerChat); ok {
		log.Println("Послал id: PeerChat")
	} else if _, ok := m.FromID.(*tg.PeerChannel); ok {
		log.Println("Послал id: PeerChannel")
	}
	log.Println(m.FromID) // почемуто nil
	log.Println("------------0000000111111111100000-------------")
	// Peer ID, чат, в который было отправлено это сообщение
	if peerUser, ok := m.PeerID.(*tg.PeerUser); ok {
		log.Println("Послал id: ", peerUser.UserID)
	} else if _, ok := m.PeerID.(*tg.PeerChat); ok {
		log.Println("Послал id: PeerChat")
	} else if _, ok := m.PeerID.(*tg.PeerChannel); ok {
		log.Println("Послал id: PeerChannel")
	}
	log.Println(m.FromID) // почемуто nil
	//! https://core.telegram.org/constructor/message
	log.Println("ЖЖ>", m.Message)

	log.Println("------------1111111111111111-------------")
	//	log.Println(utilites.FormatObject(u))
	// смотрим чат в который было отправлено сообщение .. хз
	if _, ok := m.PeerID.(*tg.PeerChat); ok { // Peer ID, чат, в который было отправлено это сообщение
		log.Println("Это чат")
	} else if _, ok := m.PeerID.(*tg.PeerChannel); ok {
		log.Println("Это канал")
	} else if peerUser, ok := m.PeerID.(*tg.PeerUser); ok {
		log.Println("Это пользователь")
		userid := peerUser.UserID //
		log.Println("User Id:", userid)
		log.Println("ID message:", m.ID)
		d := entities.Users[userid]
		log.Println("user", d)
		log.Println("------------22222222222222222----------")
		// entities.Users[userid]
		//! https://core.telegram.org/constructor/user
		for key, user := range entities.Users { // в key ID user
			if user.Self {
				log.Println("Этот пользователь: ", user.ID, " FirstName: ", user.FirstName, "Бот? ", user.Bot, user.AccessHash)
			} else {
				log.Println("Пользователь: ", user.ID, " FirstName: ", user.FirstName, "Бот? ", user.Bot, user.AccessHash)
			}
			if key == m.ID {
				log.Println("отправил: ", user.ID, " FirstName: ", user.FirstName, "Бот? ", user.Bot, user.AccessHash)
			} else {
				log.Println(key, m.ID)
			}
		}
		log.Println("------------44444444444444444444444444444444----------")

	}

	// Отправка ответа.
	_, err := sender.Reply(entities, u).Text(ctx, "Я получил:"+m.Message)
	// rr := sender.To(&tg.InputPeerChat{ChatID: 22222222})
	// rr.Builder.Text(ctx, "Хай")

	return err
}
func onNewMessage(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
	//! https://core.telegram.org/constructor/updateNewMessage
	m, ok := u.Message.(*tg.Message) // *tg.MessageService
	if !ok || m.Out {
		//Исходящее сообщение, не интересно.
		return nil
	}
	//! https://core.telegram.org/constructor/user
	if peerUser, ok := m.PeerID.(*tg.PeerUser); ok {
		userid := peerUser.UserID
		log.Println("Послал id: ", userid)
		if user, ok := entities.Users[userid]; ok {
			log.Println("User:", user.ID, " FirstName: ", user.FirstName, "Бот? ", user.Bot, user.AccessHash)
		}
	}

	return nil
}
