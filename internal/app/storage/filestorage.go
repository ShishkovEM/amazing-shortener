package storage

import (
	"encoding/json"
	"os"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteLink(link *Link) error {
	return p.encoder.Encode(&link)
}

func (p *producer) Close() error {
	return p.file.Close()
}
