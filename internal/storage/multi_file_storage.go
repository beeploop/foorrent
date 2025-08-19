package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/utils"
)

type MultiFileStorage struct {
	outDir      string
	filemap     []metadata.FileEntry
	openedFiles []*os.File
}

func NewMultiFileStorage(outDir string, filemap []metadata.FileEntry) *MultiFileStorage {
	return &MultiFileStorage{
		outDir:      outDir,
		filemap:     filemap,
		openedFiles: make([]*os.File, len(filemap)),
	}
}

func (s *MultiFileStorage) Init() error {
	for i, entry := range s.filemap {
		fullpath := filepath.Join(s.outDir, entry.Path)
		if err := os.MkdirAll(filepath.Dir(fullpath), 0777); err != nil {
			return err
		}

		file, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}

		s.openedFiles[i] = file
	}

	return nil
}

func (s *MultiFileStorage) WritePiece(index, length int, data []byte) error {
	pieceStart := int64(index) * int64(length)
	pieceEnd := pieceStart + int64(len(data))

	targetFileIndexes := s.targetFileIndex(pieceStart, pieceEnd)
	if utils.IsEmpty(targetFileIndexes) {
		err := fmt.Errorf("piece does not map to any opened files")
		return err
	}

	// TODO: Write the piece to the target files
	// a piece can span multiple files, handle it

	return nil
}

func (s *MultiFileStorage) Close() error {
	for _, file := range s.openedFiles {
		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Returns the index of the files in s.openedFiles where the piece overlaps
func (s *MultiFileStorage) targetFileIndex(pieceStart, pieceEnd int64) []int {
	fileIndexes := make([]int, 0)
	for i, entry := range s.filemap {
		pieceRange := utils.Range{
			Start: pieceStart,
			End:   pieceEnd,
		}

		entryRange := utils.Range{
			Start: entry.Begin(),
			End:   entry.End(),
		}

		if utils.IsOverlapping(pieceRange, entryRange) {
			fileIndexes = append(fileIndexes, i)
		}
	}
	return fileIndexes
}
