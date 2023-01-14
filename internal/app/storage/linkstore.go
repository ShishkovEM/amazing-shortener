package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/speps/go-hashids"
)

// Link Структура записи информации о гиперссылках
type Link struct {
	ID       int    `json:"id"`       // Идентификатор гиперссылки
	Original string `json:"original"` // Исходная (длинная) ссылка
	Short    string `json:"short"`    // Короткая ссылка
}

// LinkStore Структура для хранения записей типа Link в оперативной памяти
type LinkStore struct {
	sync.Mutex

	fileName string
	links    map[string]Link
	nextID   int
}

// NewLinkStore Создаёт новый LinkStore
func NewLinkStore(fileName string) *LinkStore {
	ls := &LinkStore{}
	ls.links = make(map[string]Link)
	ls.nextID = 0
	ls.fileName = fileName

	if ls.fileName != "" {
		file, err := os.Open(fileName)

		if strings.Contains(err.Error(), "no such file or directory") {
			file, err = os.Create(fileName)
			if err != nil {
				log.Fatalf("Error when creating file: %s", err)
			}
		} else if err != nil {
			log.Fatalf("Error when opening file: %s", err)
		}
		fileScanner := bufio.NewScanner(file)
		lineCounter := 1
		for fileScanner.Scan() {
			link := Link{}
			err := json.Unmarshal(fileScanner.Bytes(), &link)
			if err != nil {
				log.Fatalf("Error when unmarshalling file %s at line %d", err, lineCounter)
			}
			ls.addLinkToMemStorage(link)
			lineCounter++
		}
		if err := fileScanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
		}
		err = file.Close()
		if err != nil {
			log.Fatalf("Error when closing file: %s", err)
		}
	}
	return ls
}

func (ls *LinkStore) addLinkToMemStorage(link Link) {
	ls.Lock()
	defer ls.Unlock()

	link.ID = ls.nextID
	ls.links[link.Short] = link
	ls.nextID++
}

// CreateLink создаёт новую запись в LinkStore
func (ls *LinkStore) CreateLink(longURL string) string {
	ls.Lock()
	defer ls.Unlock()

	link := Link{
		ID:       ls.nextID,
		Original: longURL,
		Short:    shorten(),
	}

	if ls.fileName == "" {
		ls.links[link.Short] = link
		ls.nextID++
	} else {
		ls.links[link.Short] = link
		ls.nextID++
		producer, err := newProducer(ls.fileName)
		if err != nil {
			log.Fatal(err)
		}
		err = producer.WriteLink(&link)
		if err != nil {
			log.Fatal(err)
		}
		err = producer.Close()
		if err != nil {
			log.Fatalf("Error when closing producer: %s", err)
		}
	}
	return link.Short
}

// shorten() Создаёт короткий хэш
func shorten() string {
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	short, _ := h.Encode([]int{int(now.UnixMicro())})
	return short
}

// GetLink получает запись об одной ссылке по её id
func (ls *LinkStore) GetLink(short string) (Link, error) {
	ls.Lock()
	defer ls.Unlock()

	l, ok := ls.links[short]

	if ok {
		return l, nil
	} else {
		return Link{}, fmt.Errorf("link with id=%s not found", short)
	}
}
