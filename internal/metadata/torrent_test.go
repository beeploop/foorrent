package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTorrentMetadata(t *testing.T) {
	t.Run("Test Filemap func", func(t *testing.T) {
		torrent := Torrent{
			Info: TorrentInfo{
				Files: []TorrentFile{
					{
						Path:   []string{"CorePlus-current.iso"},
						Length: 124780544,
					},
					{
						Path:   []string{"TCL_CorePlusCurrent_meta.sqlite"},
						Length: 9216,
					},
					{
						Path:   []string{"TCL_CorePlusCurrent_meta.xml"},
						Length: 1013,
					},
				},
			},
		}

		filemap := torrent.FileMap()

		assert.Equal(t, int64(0), filemap[0].Offset)
		assert.Equal(t, int64(124780544), filemap[1].Offset)
		assert.Equal(t, int64(124789760), filemap[2].Offset)
	})
}
