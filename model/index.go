package model

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
)

const (
	notesPerPage = 10
	pagesBucket  = "pages"
)

// RebuildIndex completeky rebuilds index from scratch
func RebuildIndex(db *DB) error {
	var timestamps []int
	notesMap := make(map[int]string)

	if err := db.View(func(tx *bolt.Tx) error {
		bucketNotes := tx.Bucket([]byte(NotesBucket))
		if bucketNotes != nil {
			cursor := bucketNotes.Cursor()
			for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
				note := new(Note)
				json.Unmarshal(v, &note)

				timestamp := int(note.CreatedAt.Unix())
				if _, exists := notesMap[timestamp]; exists {
					return fmt.Errorf("Duplicate note found at %v (%s)", note.CreatedAt, note.UUID)
				}

				timestamps = append(timestamps, timestamp)
				notesMap[timestamp] = note.UUID
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if len(timestamps) == 0 {
		return nil
	}

	sort.Sort(sort.Reverse(sort.IntSlice(timestamps)))

	if err := db.Update(func(tx *bolt.Tx) error {
		pagesMap := make(map[string][]string)
		// tags := make(map[string][]*Note)

		page := 0
		pageKey := fmt.Sprintf("page-%d", page)

		bucketNotes := tx.Bucket([]byte(NotesBucket))
		if bucketNotes != nil {
			for _, timestamp := range timestamps {
				bytes := bucketNotes.Get([]byte(notesMap[timestamp]))
				note := new(Note)

				json.Unmarshal(bytes, &note)

				pagesMap[pageKey] = append(pagesMap[pageKey], note.UUID)
				if len(pagesMap[pageKey]) >= notesPerPage {
					page++
					pageKey = fmt.Sprintf("page-%d", page)
				}
			}
		}

		tx.DeleteBucket([]byte(pagesBucket))
		bucketPages, err := tx.CreateBucketIfNotExists([]byte(pagesBucket))
		if err != nil {
			return err
		}

		for k, v := range pagesMap {
			bytes, _ := json.Marshal(v)
			if err := bucketPages.Put([]byte(k), []byte(bytes)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetNotes returns list of notes fot the given page number
func GetNotes(page int, db *DB) ([]*Note, error) {
	key := fmt.Sprintf("page-%d", page)
	var uuids []string
	var notes []*Note

	if err := db.View(func(tx *bolt.Tx) error {
		bucketPages := tx.Bucket([]byte(pagesBucket))
		if bucketPages != nil {
			bytes := bucketPages.Get([]byte(key))
			json.Unmarshal(bytes, &uuids)

			if len(uuids) > 0 {
				bucketNotes := tx.Bucket([]byte(NotesBucket))
				if bucketNotes == nil {
					return fmt.Errorf("Unable to get Notes bucket!")
				}

				for _, uuid := range uuids {
					note := new(Note)
					uuidBytes := bucketNotes.Get([]byte(uuid))
					json.Unmarshal(uuidBytes, &note)
					if note.CreatedAt.IsZero() {
						return fmt.Errorf("Error unmarshalling note!")
					}

					notes = append(notes, note)
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return notes, nil
}
