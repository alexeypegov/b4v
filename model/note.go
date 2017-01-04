package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	// NotesBucket contains notes bucket name
	NotesBucket = "notes"
)

const (
	// Draft determines whatever this note is published or not
	Draft byte = 1 << iota

	// PlainHTML format of the entry
	PlainHTML
)

var (
	uuidRegexp = regexp.MustCompile("([\\s\\pP\\pS]+)")
)

// Note is a note, yeah
type Note struct {
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	Flags     byte      `json:"flags"`
}

// SaveAll save all the given Notes
func SaveAll(notes []Note, db *DB) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		for _, n := range notes {
			if err := saveInt(&n, false, tx); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetNote loads Note by its uuid
func GetNote(uuid string, db *DB) (*Note, error) {
	bytes, err := db.Get(NotesBucket, uuid)
	if err != nil {
		return nil, err
	}

	result := new(Note)
	json.Unmarshal(bytes, &result)
	return result, nil
}

func genUUID(title string) string {
	time := time.Now()

	replaced := strings.ToLower(uuidRegexp.ReplaceAllLiteralString(title, "-"))
	return fmt.Sprintf("%d%02d%02d-%s", time.Year(), time.Month(), time.Day(), replaced)
}

func saveInt(note *Note, draft bool, tx *bolt.Tx) error {
	bucketNotes, err := tx.CreateBucketIfNotExists([]byte(NotesBucket))
	if err != nil {
		return err
	}

	if len(note.UUID) == 0 {
		note.UUID = genUUID(note.Title)
	}

	if draft {
		note.Flags |= Draft
	}

	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now()
	}

	jsonNote, _ := json.Marshal(note)
	if err := bucketNotes.Put([]byte(note.UUID), []byte(jsonNote)); err != nil {
		return err
	}

	return nil
}

// Save Note to an underlying storage
func (note *Note) Save(draft bool, db *DB) error {
	if err := db.DB.Update(func(tx *bolt.Tx) error {
		return saveInt(note, draft, tx)
	}); err != nil {
		return err
	}

	return nil
}
