package templates

import (
	"fmt"
	"bytes"
	"testing"
	"time"
  "strconv"
	"runtime"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alexeypegov/b4v/model"
)

func getTestData(ext string) string {
	pc, path, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	funcName := parts[len(parts) - 1][4:] // skip Test prefix
	filename := fmt.Sprintf("%s/test_data/%s.%s", filepath.Dir(path), funcName, ext)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Sprintf("File was not found: '%s'", filename)
	}

	return string(data)
}

func assertEquals(s string, t *testing.T) {
	expected := getTestData("html")
	if expected != s {
		t.Errorf("\nE: %s\nA: %s", strconv.Quote(expected), strconv.Quote(s))
	}
}

func TestRenderNote(t *testing.T) {
	tpl := New(".")

	ts, _ := time.Parse(time.RFC822, "04 Nov 79 22:23 MSK")
	note := &model.Note{Title: "first", Content: "<h1>first content</h1>", Tags: []string{"саптрю", "слушаю"}, CreatedAt: ts}

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, "note", note); err != nil {
		t.Error(err)
	}

	assertEquals(w.String(), t)
}

func TestRenderNotes(t *testing.T) {
	tpl := New(".")
  ts, _ := time.Parse(time.RFC822, "04 Nov 79 22:23 MSK")

	notes := []*model.Note{}

	notes = append(notes, &model.Note{Title: "first", Content: "first content", CreatedAt: ts})
	notes = append(notes, &model.Note{Title: "second", Content: "second content", Tags: []string{"саптрю", "слушаю"}, CreatedAt: ts})

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, "notes", notes); err != nil {
		t.Error(err)
	}

	assertEquals(w.String(), t)
}
