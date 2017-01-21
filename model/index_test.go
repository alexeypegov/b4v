package model

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFailOnDuplicateDates(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	now := time.Now()

	note1 := Note{Title: "One", CreatedAt: now}
	note2 := Note{Title: "Two", CreatedAt: now}

	note1.Save(false, db)
	note2.Save(false, db)

	if err := RebuildIndex(db); err == nil || !strings.HasPrefix(err.Error(), "Duplicate note found") {
		t.Errorf("Should fail on duplicate dates! %s", err.Error())
	}
}

func TestBuildIndex(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	if notes, err := GetNotes(0, db); notes != nil || err != nil {
		t.Error("There should be no notes yet")
	}

	for i := 0; i < 15; i++ {
		at, _ := time.Parse(time.RFC822, fmt.Sprintf("11 Nov 79 %02d:23 MSK", i))
		(&Note{Title: fmt.Sprintf("Title %d", i), CreatedAt: at}).Save(false, db)
	}

	if err := RebuildIndex(db); err != nil {
		t.Error(err)
	}

	after, err := GetNotes(0, db)
	if err != nil {
		t.Error(err)
	}

	if len(after) < 10 {
		t.Errorf("Error fetching notes for page %d (len() = %d)!", 0, len(after))
	}
}

func TestPagesCount(t *testing.T) {
	db := NewTestDB()
	defer CloseAndDestroy(db)

	if err := (&Note{Title: "One", CreatedAt: time.Now().AddDate(0, 0, 1)}).Save(false, db); err != nil {
			t.Error("Unable to save note One")
	}

	if err := (&Note{Title: "Two", CreatedAt: time.Now()}).Save(false, db); err != nil {
		t.Error("Unable to save note Two")
	}

	if err := RebuildIndex(db); err != nil {
		t.Error(err)
	}

	count, error := GetPagesCount(db)
	if error != nil {
		t.Error("Error getting pages count")
	}

	if count != 1 {
		t.Errorf("Should be exactly 1 page, got %d instead!", count)
	}
}
