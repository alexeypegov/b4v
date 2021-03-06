package model

import (
	"encoding/json"
	"encoding/xml"
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

type rssGUID struct {
	XMLName     xml.Name `xml:"guid"`
	IsPermaLink bool     `xml:"isPermaLink,attr"`
	Data        string   `xml:",chardata"`
}

type rssDescription struct {
	XMLName xml.Name `xml:"description"`
	Data    string   `xml:",cdata"`
}

// RssItem one single RSS item
type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	GUID        rssGUID
	Title       string         `xml:"title"`
	Category    string         `xml:"category"`
	PubDate     string         `xml:"pubDate"`
	Description rssDescription `xml:"description"`
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

	if bytes == nil {
		return nil, fmt.Errorf("Note not found '%s'", uuid)
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

// ToRSS convert note to RSS item
func (note *Note) ToRSS(baseURL string) RssItem {
	return RssItem{
		GUID: rssGUID{
			IsPermaLink: true,
			Data:        fmt.Sprintf("%s%s", baseURL, note.UUID),
		},
		Title:       note.Title,
		Category:    strings.Join(note.Tags, ","),
		PubDate:     note.CreatedAt.Format(time.RFC1123Z),
		Description: rssDescription{
			Data: note.Content,
		},
	}
}
