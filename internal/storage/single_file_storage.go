package storage

import (
	"fmt"
	"os"
)

type SingleFileStorage struct {
	filename string
	file     *os.File
}

func NewSingleFileStorage(filename string) *SingleFileStorage {
	return &SingleFileStorage{
		filename: filename,
	}
}

func (s *SingleFileStorage) Init() error {
	file, err := os.OpenFile(s.filename, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}

	s.file = file

	return nil
}

func (s *SingleFileStorage) WritePiece(index, length int, data []byte) error {
	offset := int64(index * length)
	if _, err := s.file.WriteAt(data, offset); err != nil {
		err := fmt.Errorf("Error occurred while writing to file, aborted")
		return err
	}

	return nil
}

func (s *SingleFileStorage) Close() error {
	return s.file.Close()
}
