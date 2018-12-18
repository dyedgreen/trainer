// Implement persist functionality for
// flashcard boxes

package flashcard

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

var (
	BoxDbBucket = "boxes"
)

type Shelve struct {
	db *bolt.DB
}

func OpenShelve(dbFile string) *Shelve {
	var s Shelve
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		panic("Could not open flashcard database")
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BoxDbBucket))
		return err
	})
	if err != nil {
		a.db.Close()
		panic("Could not create boxes bucket")
	}
	s.db = db
	return &s
}

func (s *Shelve) Close() {
	s.db.Close()
}