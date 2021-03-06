package model

import (
	"encoding/json"
	"testing"
	"time"

	"reflect"
)

func TestUnmarshalling(t *testing.T) {
	data := []byte(`[
{
	"title": "GoGo Penguin Live from Old Grenada Studios",
	"date" : "29 Ноября 2016",
	"url"  : "/note/%D0%BF%D0%B5%D1%80%D0%B2%D1%8B%D0%B9-%D0%BA%D0%BB%D0%B0%D1%81%D1%81",
	"tags" : ["слушаю"],
	"content": "<div></div>"
}]`)

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

	if "первый-класс" != note.UUID {
		t.Error("Error unmarshalling UUID", note.UUID)
	}
}

func TestIncrementTimeForDuplicateEntries(t *testing.T) {
	data := []byte(`[
{"title": "One",
 "date" : "29 Ноября 2016",
 "url"  : "/note/one",
 "tags" : ["слушаю"],
 "content": "<div></div>"
},
{"title": "Two",
 "date" : "29 Ноября 2016",
 "url"  : "/note/two",
 "tags" : ["слушаю"],
 "content": "<div></div>"
}]`)

	var notes []oldNote
	if err := json.Unmarshal(data, &notes); err != nil {
		t.Error(err)
	}

	if len(notes) != 2 {
		t.Error("Should unmarshal exactly two enties")
	}

	converted := ConvertNotes(notes)

	if converted[0].CreatedAt.Equal(converted[1].CreatedAt) {
		t.Errorf("Should increment duplicate time (%v and %v)", converted[0].CreatedAt, converted[1].CreatedAt)
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

	if result.Flags&PlainHTML != PlainHTML {
		t.Error("Old entries should be migrated as PlainHTML")
	}
}
