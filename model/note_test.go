package model

import (
  "fmt"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
  testdb := OpenTestDB()
  defer testdb.CloseTestDB()

  note := Note{
    UUID: "something-interesting",
    Title: "title",
    Content: "Some content"}

  if err := note.Save(true, testdb.DB); err != nil {
    t.Error("Unable to save note:", err)
  }

  var loaded Note
  loaded.Load("something-interesting", testdb.DB)

  if loaded.Title != "title" {
    t.Error(fmt.Sscanf("Unable to load note (wrong title: '%s' vs '%s')", "title", loaded.Title))
  }
}
