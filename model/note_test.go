package model

import (
	"fmt"
	"testing"
	"time"
)

func TestSaveAndLoad(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	note := Note{UUID: "something-interesting", Title: "title", Content: "Some content"}

	if err := note.Save(true, db); err != nil {
		t.Error("Unable to save note:", err)
	}

	loaded, err := GetNote("something-interesting", db)
	if err != nil {
		t.Error(err)
	}

	if loaded.Title != "title" {
		t.Errorf("Unable to load note (wrong title: '%s' vs '%s')", "title", loaded.Title)
	}
}

func TestNotExistingNote(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	note := Note{UUID: "abc"}
	if err := note.Save(false, db); err != nil {
		t.Error("Unable to save note:", err)
	}

	loaded, err := GetNote("not-existing", db)
	if err != nil {
		t.Error("Error loading Note")
	}

	if len(loaded.UUID) > 0 {
		t.Error("Inexisting note was found")
	}
}

func TestGenerateUUID(t *testing.T) {
  now := time.Now()
	prefix := fmt.Sprintf("%4d%2d%2d", now.Year(), now.Month(), now.Day())

  s := genUUID("hello");
  if s != prefix + "-hello" {
    t.Errorf("Invalid uuid generated: '%s'", s)
  }

  s = genUUID("ЮНИКОД и 123")
  if s != prefix + "-юникод-и-123" {
    t.Errorf("Invalid uuid generated: '%s'", s)
  }

  s = genUUID("разное &^%$#@!*(){}[]?/\\ такое")
  if s != prefix + "-разное-такое" {
    t.Errorf("Invalid uuid generated: '%s'", s)
  }
}

func TestAssingUUID(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

  now := time.Now()
	uuid := fmt.Sprintf("%d%d%d-а-тут-у-нас-будет-какой-то-такой-заголовок", now.Year(), now.Month(), now.Day())

	note := Note{Title: "А тут у нас будет какой-то такой заголовок", Content: "nope"}
	if err := note.Save(false, db); err != nil {
		t.Error("Unable to save note:", err)
	}

  if uuid != note.UUID {
		t.Errorf("UUID should be assigned, found: '%s' (should be '%s')", note.UUID, uuid)
	}

	loaded, err := GetNote(uuid, db)
	if err != nil {
		t.Error("Unable to load note by newly getnerated UUID")
	}

	if uuid != loaded.UUID {
		t.Errorf("UUID should be assigned, found: '%s'", loaded.UUID)
	}

  if "nope" != loaded.Content {
		t.Errorf("Content is different from the expected one, found: '%s'", loaded.Content)
	}
}
