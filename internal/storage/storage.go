package storage

import "fmt"

type Storage interface {
	WriteBlock(index int, offset int, data []byte) error
	Downloaded() int
}

type MockStorage struct{}

func (s *MockStorage) WriteBlock(index, offset int, data []byte) error {
	fmt.Printf("written block index: %d, offset: %d\n", index, offset)
	return nil
}

func (s *MockStorage) Downloaded() int {
	return 0
}
