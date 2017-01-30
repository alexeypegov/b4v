package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alexeypegov/b4v/model"
)

func getTestDataPath(ext string) (string, bool) {
	i := 1
	for ; ; i++ {
		pc, path, _, ok := runtime.Caller(i)
		if !ok {
			break
		}

		parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		last := parts[len(parts)-1]
		if strings.HasPrefix(last, "Test") {
			funcName := last[4:] // skip Test prefix
			filename := fmt.Sprintf("%s/test_data/%s.%s", filepath.Dir(path), funcName, ext)
			return filename, true
		}
	}

	return "", false
}

func getTestData(ext string) string {
	path, ok := getTestDataPath(ext)
	if !ok {
		return fmt.Sprint("Should be called from a test function!")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("File was not found: '%s'", path)
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
	note := &model.Note{
		UUID:      "first",
		Title:     "first",
		Content:   "<h1>first content</h1>",
		Tags:      []string{"саптрю", "слушаю"},
		CreatedAt: ts}

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, "note", note); err != nil {
		t.Fatal(err)
	}

	assertEquals(w.String(), t)
}

func TestRenderNotes(t *testing.T) {
	tpl := New(".")
	ts, _ := time.Parse(time.RFC822, "04 Nov 79 22:23 MSK")

	notes := []*model.Note{}

	notes = append(notes, &model.Note{
		UUID:      "first",
		Title:     "first",
		Content:   "first content",
		CreatedAt: ts})

	notes = append(notes, &model.Note{
		UUID:      "second",
		Title:     "second",
		Content:   "second content",
		Tags:      []string{"саптрю", "слушаю"},
		CreatedAt: ts})

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, "notes", notes); err != nil {
		t.Fatal(err)
	}

	assertEquals(w.String(), t)
}

func paginationTest(template string, page, total int, t *testing.T) {
	tpl := New(".")

	data := map[string]interface{}{
		"Vars": map[string]string{
			"PreviousPage": "prev",
			"NextPage":     "next",
		},
		"Paging": map[string]int{
			"Current": page,
			"Total":   total,
		},
	}

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, template, data); err != nil {
		t.Fatal(err)
	}

	assertEquals(w.String(), t)
}

func TestPaginationStart1(t *testing.T) {
	paginationTest("paging-start", 1, 1, t)
}

func TestPaginationStart2(t *testing.T) {
	paginationTest("paging-start", 2, 2, t)
}

func TestPaginationStart3(t *testing.T) {
	paginationTest("paging-start", 3, 3, t)
}

func TestPaginationEnd1(t *testing.T) {
	paginationTest("paging-end", 1, 1, t)
}

func TestPaginationEnd2(t *testing.T) {
	paginationTest("paging-end", 1, 2, t)
}

func TestFullPage(t *testing.T) {
	tpl := New(".")

	ts, _ := time.Parse(time.RFC822, "04 Nov 79 22:23 MSK")

	notes := []*model.Note{}

	notes = append(notes, &model.Note{
		UUID:      "first",
		Title:     "first",
		Content:   "first content",
		CreatedAt: ts})

	notes = append(notes, &model.Note{
		UUID:      "second",
		Title:     "second",
		Content:   "second content",
		Tags:      []string{"саптрю", "слушаю"},
		CreatedAt: ts.Add(10 * time.Minute)})

	data := map[string]interface{}{
		"Notes": notes,
		"Vars": map[string]string{
			"PreviousPage": "prev",
			"NextPage":     "next",
		},
		"Paging": map[string]int{
			"Current": 1,
			"Total":   2,
		},
	}

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, "index.tpl", data); err != nil {
		t.Fatal(err)
	}

	assertEquals(w.String(), t)
}
