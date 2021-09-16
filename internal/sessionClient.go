package internal

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ark-go/codefindbot/internal/cripto"
	"github.com/ark-go/codefindbot/internal/jt"
	"github.com/gotd/td/session"
)

// memorySession implements in-memory session storage.
// Goroutine-safe.
type memorySession struct {
	mux      sync.RWMutex
	data     []byte
	fileName string
}

// Использование пользовательского хранилища сеансов.
// Вы можете сохранить сеанс в базе данных, например Redis, MongoDB или postgres.
// Подробности реализации смотрите в memorySession.
var sessionStorage *memorySession

func init() {
	//	sessionStorage = &memorySession{}
}
func (s *memorySession) New(fname string) *memorySession {
	ms := &memorySession{}
	ms.fileName = fname
	return ms
}

// LoadSession loads session from memory.
func (s *memorySession) LoadSession(context.Context) ([]byte, error) {
	log.Println("Чтение сессии !!!!")
	if s == nil {
		return nil, session.ErrNotFound
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	if len(s.data) == 0 {
		s.ReadSession()
	}
	if len(s.data) == 0 {
		return nil, session.ErrNotFound
	}

	cpy := append([]byte(nil), s.data...)

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *memorySession) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	s.data = data
	s.SaveSession()
	s.mux.Unlock()
	return nil
}

func (s *memorySession) SaveSession() {
	permissions := 0644 // or whatever you need
	filename := filepath.Join(jt.RootDir, s.fileName)
	criptdata, err := cripto.EncriptFile(s.data)
	if err != nil {
		log.Println("Ошибка шифрования сессии", err)
		return
	}
	if err := os.WriteFile(filename, criptdata, fs.FileMode(permissions)); err != nil {
		log.Println("Ошибка записи сессии на диск")
	} else {
		log.Println("Запись сессии на диск")
	}
}
func (s *memorySession) ReadSession() {
	filename := filepath.Join(jt.RootDir, s.fileName)
	criptdata, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Ошибка чтения сессии c диска", err)
		return
	} else {
		log.Println("Читаем сессию из файла", filename)
	}
	if decriptdata, err := cripto.DecriptFile(criptdata); err != nil {
		log.Println("Ошибка расшифровки сессии с диска", err)
	} else {
		s.data = decriptdata
	}
}
