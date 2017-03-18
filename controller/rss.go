package controller

import (
	"encoding/xml"
	"net/http"

	"github.com/alexeypegov/b4v/model"
)

// Rss20 Rss 2.0 main info
type Rss20 struct {
	XMLName     xml.Name        `xml:"rss"`
  Version     string          `xml:"version,attr"`
	Title       string          `xml:"channel>title"`
	Link        string          `xml:"channel>link"`
	Description string          `xml:"channel>description"`
	TTL         int             `xml:"channel>ttl"`
	Items       []model.RssItem `xml:"channel>item"`
	// Image string `xml:"channel>image"`
}

func convert(baseURL string, notes []*model.Note) []model.RssItem {
	result := make([]model.RssItem, len(notes))
	for i, n := range notes {
		result[i] = n.ToRSS(baseURL)
	}

	return result
}

func generateRss(notes []*model.Note, title, link, desc string) ([]byte, error) {
  ba, err := xml.MarshalIndent(Rss20{
    Version: "2.0",
		Title:       title,
		Link:        link,
		Description: desc,
		TTL:         60,
		Items:       convert(link, notes),
	}, "", "  ")
  
	if err != nil {
		return nil, err
	}
  
  return ba, nil
}

// RssHandler handles RSS requests
func RssHandler(w http.ResponseWriter, r *http.Request, ctx *Context) (int, error) {
	notes, err := model.GetNotes(1, ctx.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	ba, err := generateRss(notes, ctx.Vars["title"], ctx.Vars["url"], ctx.Vars["title"])
  if err != nil {
		return http.StatusInternalServerError, err
	}
  
  w.Write([]byte(xml.Header))
  w.Write(ba)

	return http.StatusOK, nil
}
