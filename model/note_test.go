package model

import (
  "fmt"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
  testdb := OpenTestDB()
  defer testdb.CloseTestDB()

  note := Note{UUID: "something-interesting",Title: "title",Content: "Some content"}

  if err := note.Save(true, testdb.DB); err != nil {
    t.Error("Unable to save note:", err)
  }

  loaded, err := GetNote("something-interesting", testdb.DB)
  if err != nil {
    t.Error(err)
  }
  
  if loaded.Title != "title" {
    t.Error(fmt.Sscanf("Unable to load note (wrong title: '%s' vs '%s')", "title", loaded.Title))
  }
}

func TestGenerateUUID(t *testing.T) {
  testdb := OpenTestDB()
  defer testdb.Close()

  note := Note{Title: "А тут у нас будет какой-то такой заголовок", Content: "nope"}

  if err := note.Save(false, testdb.DB); err != nil {
    t.Error("Unable to save note:", err)
  }


}
