package storage

import "fmt"

type Storage interface {
	Init() error
	WritePiece(index int, length int, data []byte) error
	Close() error
}

type MockStorage struct{}

func (s *MockStorage) Init() error {
	return nil
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (s *MockStorage) WritePiece(index, length int, data []byte) error {
	fmt.Printf("written block index: %d, offset: %d\n", index, length)
	return nil
}

func (s *MockStorage) Close() error {
	return nil
}
