package model

import (
	"encoding/json"
	"testing"
	"time"

	"reflect"
)

func TestUnmarshalling(t *testing.T) {
	data := []byte("[{\"title\": \"GoGo Penguin Live from Old Grenada Studios\",\"date\": \"29 Ноября 2016\",\"url\": \"/note/gogo-penguin-live-from-old-grenada-studios\",\"tags\": [\"слушаю\"],\"content\": \"<div></div>\"}]")

	var notes []oldNote
  if err := json.Unmarshal(data, &notes); err != nil {
    t.Error(err)
  }

  if len(notes) != 1 {
    t.Error("Should unmarshal exactly one entry")
  }

  note := notes[0]
  if note.Title != "GoGo Penguin Live from Old Grenada Studios" {
    t.Error("Error unmarshalling title")
  }

  formattedDate := note.Date.Format(time.RFC822)
  if "29 Nov 16 12:00 MSK" != formattedDate {
    t.Error("Error unmarshalling date", formattedDate)
  }

  if note.Content != "<div></div>" {
    t.Error("Error unmarshalling content")
  }

  if !reflect.DeepEqual(note.Tags, []string{"слушаю"}) {
    t.Error("Error unmarshalling tags", note.Tags)
  }

  if "gogo-penguin-live-from-old-grenada-studios" != note.UUID {
    t.Error("Error unmarshalling UUID", note.UUID)
  }
}

func TestToNote(t *testing.T) {
	n := oldNote{UUID: "haha", Title: "title", Date: time.Now(), Tags: []string{"a"}, Content: "content"}
	result := n.toNote()
	if result.UUID != n.UUID {
		t.Error("UUID should be migrated")
	}

	if result.Title != n.Title {
		t.Error("Title should be migrated")
	}

	if result.CreatedAt != n.Date {
		t.Error("Date should be migrated")
	}

	if result.Content != n.Content {
		t.Error("Content should be migrated")
	}

	if !reflect.DeepEqual(result.Tags, n.Tags) {
		t.Error("Tags should be migrated")
	}

  if result.Flags & PlainHTML != PlainHTML {
    t.Error("Old entries should be migrated as PlainHTML")
  }
}
