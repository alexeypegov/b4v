package templates

import (
	"bytes"
	"testing"
	"time"
  "strconv"

	"github.com/alexeypegov/b4v/model"
)

func TestRenderNote(t *testing.T) {
	tpl := New(".")

	ts, _ := time.Parse(time.RFC822, "11 Nov 79 22:23 MSK")
	note := &model.Note{Title: "first", Content: "<h1>first content</h1>", Tags: []string{"саптрю", "слушаю"}, CreatedAt: ts}

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, "index.tpl", Data{Note: note}); err != nil {
		t.Error(err)
	}

	body := `<!doctype html>
<html>
<head>
  <title></title>
</head>
<body>
<div class="note">
  <div class="title">first</div>
  <div class="date">1979-11-11 22:23:00 &#43;0000 MSK</div>
  <div class="tags"><div class="tag">саптрю</div><div class="tag">слушаю</div></div>
  <div class="body"><h1>first content</h1></div>
</div>
</body>
</html>
`

	out := w.String()
	if out != body {
		t.Errorf("\nE: %s\nA: %s", strconv.Quote(body), strconv.Quote(out))
	}
}

func TestRenderNotes(t *testing.T) {
	tpl := New(".")
  ts, _ := time.Parse(time.RFC822, "11 Nov 79 22:23 MSK")

	notes := []*model.Note{}

	notes = append(notes, &model.Note{Title: "first", Content: "first content", CreatedAt: ts})
	notes = append(notes, &model.Note{Title: "second", Content: "second content", Tags: []string{"саптрю", "слушаю"}, CreatedAt: ts})

	w := bytes.NewBufferString("")
	if err := tpl.ExecuteTemplate(w, "index.tpl", Data{Notes: notes}); err != nil {
		t.Error(err)
	}

	body := `<!doctype html>
<html>
<head>
  <title></title>
</head>
<body>
<div class="note">
  <div class="title">first</div>
  <div class="date">1979-11-11 22:23:00 &#43;0000 MSK</div>
  <div class="body">first content</div>
</div>

<div class="note">
  <div class="title">second</div>
  <div class="date">1979-11-11 22:23:00 &#43;0000 MSK</div>
  <div class="tags"><div class="tag">саптрю</div><div class="tag">слушаю</div></div>
  <div class="body">second content</div>
</div>
</body>
</html>
`

  out := w.String()
	if out != body {
		t.Errorf("\nE: %s\nA: %s", strconv.Quote(body), strconv.Quote(out))
	}
}
