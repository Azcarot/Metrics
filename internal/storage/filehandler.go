package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func NewProducer(fileName string) (*Producer, error) {
	exist, err := fileOrPathExists(fileName)
	if err != nil {
		return nil, err
	}
	if !exist {
		if err := os.MkdirAll(filepath.Dir(fileName), 0777); err != nil {
			return nil, err
		}
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *Metrics) error {
	return p.encoder.Encode(&event)
}

func fileOrPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewConsumer(fileName string) (*Consumer, error) {
	exist, err := fileOrPathExists(fileName)

	if err != nil {
		return nil, err
	}
	if !exist {
		if err := os.MkdirAll(filepath.Dir(fileName), 0777); err != nil {
			return nil, err
		}
	}
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*[]Metrics, error) {
	scanner := bufio.NewScanner(c.file)
	var result []Metrics
	for scanner.Scan() {
		event := &Metrics{}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		err := json.Unmarshal([]byte(scanner.Text()), &event)
		if err != nil {
			return nil, err
		}
		result = append(result, *event)

	}
	return &result, nil
}
