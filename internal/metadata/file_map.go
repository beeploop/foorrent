package metadata

type FileEntry struct {
	Path   string
	Length int64
	Offset int64
}

func (f *FileEntry) Begin() int64 {
	return f.Offset
}

func (f *FileEntry) End() int64 {
	return f.Offset + f.Length
}
