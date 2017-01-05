package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var monthsMapping = map[string]time.Month{
	"января":   time.January,
	"февраля":  time.February,
	"марта":    time.March,
	"апреля":   time.April,
	"мая":      time.May,
	"июня":     time.June,
	"июля":     time.July,
	"августа":  time.August,
	"сентября": time.September,
	"октября":  time.October,
	"ноября":   time.November,
	"декабря":  time.December}

type oldNote struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UUID    string    `json:"url"`
	Tags    []string  `json:"tags"`
	Content string    `json:"content"`
}

// Populate import all the notes from the old format backup file
func Populate(filename string, db *DB) (int, error) {
	notes, err := ImportOldNotes(filename)
	if err != nil {
		return -1, err
	}

	if err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(NotesBucket))
		if bucket != nil {
			return fmt.Errorf("Existing Notes bucket was found")
		}

		return nil
	}); err != nil {
		return -1, err
	}

	if err := SaveAll(notes, db); err != nil {
		return -1, err
	}

	return len(notes), nil
}

// ImportOldNotes imports notes from an old format backup file
func ImportOldNotes(filename string) ([]Note, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var notes []oldNote
	if err := json.Unmarshal(content, &notes); err != nil {
		return nil, err
	}

	result := ConvertNotes(notes)
	return result, nil
}

// ConvertNotes converts slice of old notes to the new format (and fixes duplicate time)
func ConvertNotes(notes []oldNote) []Note {
	timestamps := make(map[int]bool, len(notes))

	result := make([]Note, len(notes))
	for i, n := range notes {
		converted := n.toNote()

		for {
			createdAt := converted.CreatedAt
			if _, existing := timestamps[int(createdAt.Unix())]; existing {
				converted.CreatedAt = createdAt.Add(time.Hour * 1)
			} else {

				timestamps[int(converted.CreatedAt.Unix())] = true
				break
			}
		}

		result[i] = converted
	}

	return result
}

func (n *oldNote) toNote() Note {
	return Note{UUID: n.UUID, Title: n.Title, Content: n.Content, Tags: n.Tags, CreatedAt: n.Date, Flags: PlainHTML}
}

func (n *oldNote) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		switch strings.ToLower(k) {
		case "title":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("Not a string for key[%s] %q", k, v)
			}

			n.Title = s
			break
		case "date":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("Not a string for key[%s] %q", k, v)
			}

			var day, year int
			var monthRu string
			fmt.Sscanf(s, "%d %s %d", &day, &monthRu, &year)

			month, ok := monthsMapping[strings.ToLower(monthRu)]
			if !ok {
				return fmt.Errorf("Unknown month '%s'", monthRu)
			}

			location, err := time.LoadLocation("Europe/Moscow")
			if err != nil {
				return err
			}

			n.Date = time.Date(year, month, day, 12, 0, 0, 0, location)
			break
		case "url":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("Not a string for key[%s] %q", k, v)
			}

			var uuid string
			fmt.Sscanf(s, "/note/%s", &uuid)
			n.UUID = uuid
			break
		case "tags":
			arr, ok := v.([]interface{})
			if !ok {
				return fmt.Errorf("Not an array for key[%s]: %q", k, v)
			}

			n.Tags = make([]string, len(arr))
			for i, t := range arr {
				tag, _ := t.(string)
				n.Tags[i] = tag
			}
			break
		case "content":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("Not a string for key[%s]: %q", k, v)
			}

			n.Content = s
			break
		}
	}

	return nil
}
