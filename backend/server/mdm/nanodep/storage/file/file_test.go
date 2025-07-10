package file

import (
	"testing"

	"github.com/notawar/mobius/backend/server/mdm/nanodep/storage"
	"github.com/notawar/mobius/backend/server/mdm/nanodep/storage/storagetest"
)

func TestFileStorage(t *testing.T) {
	storagetest.Run(t, func(t *testing.T) storage.AllDEPStorage {
		s, err := New(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		return s
	})
}
