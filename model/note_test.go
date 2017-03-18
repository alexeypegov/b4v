package model

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"
	
	"github.com/alexeypegov/b4v/test"
)

func TestSaveAndLoad(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	note := Note{
		UUID:    "локальзованный-урл",
		Title:   "title",
		Content: "Some content"}

	if err := note.Save(true, db); err != nil {
		t.Error("Unable to save note:", err)
	}

	loaded, err := GetNote("локальзованный-урл", db)
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

	_, err := GetNote("not-existing", db)
	if err != nil {
		if err.Error() != "Note not found 'not-existing'" {
			t.Fatal("Error loading Note:", err)
		}
	} else {
		t.Error("Inexisting note was found")
	}
}

func TestGenerateUUID(t *testing.T) {
	now := time.Now()
	prefix := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())

	s := genUUID("hello")
	if s != prefix+"-hello" {
		t.Errorf("Invalid uuid generated: '%s'", s)
	}

	s = genUUID("ЮНИКОД и 123")
	if s != prefix+"-юникод-и-123" {
		t.Errorf("Invalid uuid generated: '%s'", s)
	}

	s = genUUID("разное &^%$#@!*(){}[]?/\\ такое")
	if s != prefix+"-разное-такое" {
		t.Errorf("Invalid uuid generated: '%s'", s)
	}
}

func TestKeepOriginalCreatedAtOnSave(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	ts, _ := time.Parse(time.RFC822, "11 Nov 79 22:23 MSK")
	note := Note{Title: "Keep trying", CreatedAt: ts}
	if err := note.Save(false, db); err != nil {
		t.Error("Unable to save note", err)
	}

	if !note.CreatedAt.Equal(ts) {
		t.Errorf("CreatedAt should not be overwritten (%v vs %v)", ts, note.CreatedAt)
	}
}

func TestAssignUUID(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	now := time.Now()
	uuid := fmt.Sprintf("%d%02d%02d-а-тут-у-нас-будет-какой-то-такой-заголовок", now.Year(), now.Month(), now.Day())

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

func TestToRss(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	ts, _ := time.Parse(time.RFC822, "11 Nov 79 22:23 MSK")

	note := Note{
		UUID:      "локальзованный-урл",
		Title:     "title",
		Content:   "Some content",
		CreatedAt: ts,
		Tags:      []string{"a", "b"},
	}

	if err := note.Save(true, db); err != nil {
		t.Error("Unable to save note:", err)
	}

	rss := note.ToRSS("http://localhost/")

	ba, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		t.Fatal("Unable to marshal item to RSS")
	}

	expected := `<item>
  <guid isPermaLink="true">http://localhost/локальзованный-урл</guid>
  <title>title</title>
  <category>a,b</category>
  <pubDate>Sun, 11 Nov 1979 22:23:00 +0300</pubDate>
  <description><![CDATA[Some content]]></description>
</item>`
	
	test.AssertEquals(expected, string(ba), t)
}
