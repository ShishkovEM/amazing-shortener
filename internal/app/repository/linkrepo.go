package repository

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"
)

type LinkRepository struct {
	sync.Mutex

	InMemory   *storage.LinkStore
	repository *LinkFileRepo
	size       int
}

type LinkFileRepo struct {
	producer *Producer
	consumer *Consumer
}

type Producer struct {
	sync.Mutex

	file   *os.File
	writer *bufio.Writer
}

type Consumer struct {
	sync.Mutex

	file    *os.File
	scanner *bufio.Scanner
}

func NewLinkRepository(fileName string, linkStore *storage.LinkStore) (*LinkRepository, error) {
	producer, err := NewProducer(fileName)
	if err != nil {
		return nil, err
	}

	consumer, err := NewConsumer(fileName)
	if err != nil {
		return nil, err
	}

	lfr := LinkFileRepo{
		producer: producer,
		consumer: consumer,
	}

	linkRepository := LinkRepository{
		InMemory:   linkStore,
		repository: &lfr,
		size:       0,
	}

	if fileName != "" {
		file, err := os.OpenFile(fileName, syscall.O_RDONLY|syscall.O_CREAT, 0777)
		if err != nil {
			log.Fatalf("Error when opening/creating file: %s", err)
		}

		fileScanner := bufio.NewScanner(file)
		lineCounter := 1
		for fileScanner.Scan() {
			link := storage.Link{}
			err := json.Unmarshal(fileScanner.Bytes(), &link)
			if err != nil {
				log.Fatalf("Error when unmarshalling file %s at line %d", err, lineCounter)
			}
			linkRepository.size++
			linkRepository.InMemory.AddLinkToMemStorage(link)
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

	return &linkRepository, nil
}

func (lr *LinkRepository) Refresh(fileName string) {
	lr.Lock()
	defer lr.Unlock()

	for {
		itemsInMemStorage := lr.InMemory.GetSize()
		if lr.repository.producer.file != nil && itemsInMemStorage > lr.size {
			lr.repository.producer.Lock()

			newProducer, err := RefreshProducer(fileName)
			if err != nil {
				log.Fatalf("Error when renewing producer: %s", err)
				return
			}

			for _, link := range lr.InMemory.Links {
				err := newProducer.WriteLink(&link)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = newProducer.Close()
			if err != nil {
				log.Fatalf("Error when closing producer: %s", err)
				return
			}

			lr.repository.producer.Unlock()
		}

		time.Sleep(1 * time.Millisecond)
	}
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func RefreshProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0777)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteLink(link *storage.Link) error {
	p.Lock()
	defer p.Unlock()

	data, err := json.Marshal(&link)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *Producer) Close() error {
	p.Lock()
	defer p.Unlock()

	return p.file.Close()
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadLink() (*storage.Link, error) {
	c.Lock()
	defer c.Unlock()

	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}

	data := c.scanner.Bytes()

	link := storage.Link{}

	err := json.Unmarshal(data, &link)
	if err != nil {
		return nil, err
	}

	return &link, nil

}

func (c *Consumer) Close() error {
	c.Lock()
	defer c.Unlock()

	return c.file.Close()
}
