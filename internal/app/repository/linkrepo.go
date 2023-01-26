package repository

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
	"syscall"

	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"
)

type LinkFileRepository struct {
	sync.Mutex

	fileName   string
	repository *FileRepo
}

type FileRepo struct {
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

func NewLinkFileRepository(fileName string) (*LinkFileRepository, error) {
	if fileName != "" {
		producer, err := NewProducer(fileName)
		if err != nil {
			return nil, err
		}

		consumer, err := NewConsumer(fileName)
		if err != nil {
			return nil, err
		}

		lfr := FileRepo{
			producer: producer,
			consumer: consumer,
		}

		linkRepository := LinkFileRepository{
			fileName:   fileName,
			repository: &lfr,
		}

		return &linkRepository, nil
	}
	return nil, nil
}

func (lfr *LinkFileRepository) InitLinkStoreFromRepository(store *storage.LinkStore) {
	if lfr.fileName != "" {
		file, err := os.OpenFile(lfr.fileName, syscall.O_RDONLY|syscall.O_CREAT, 0777)
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
			store.AddLinkToMemStorage(link)
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
}

func (lfr *LinkFileRepository) WriteLinkToRepository(link *storage.Link) error {
	lfr.Lock()
	defer lfr.Unlock()

	data, err := json.Marshal(&link)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(lfr.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	if _, err := writer.Write(data); err != nil {
		return err
	}

	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
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
