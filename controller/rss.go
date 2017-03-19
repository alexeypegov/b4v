package controller

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/alexeypegov/b4v/model"
)

type rssImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

type atomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

// Rss20 Rss 2.0 main info
type Rss20 struct {
	XMLName     xml.Name        `xml:"rss"`
	Version     string          `xml:"version,attr"`
	Title       string          `xml:"channel>title"`
	Link        string          `xml:"channel>link"`
	Description string          `xml:"channel>description"`
	TTL         int             `xml:"channel>ttl"`
	Ns          string          `xml:"xmlns:atom,attr"`
	Image       rssImage        `xml:"channel>image"`
	Atom        atomLink        `xml:"channel>atom:link"`
	Items       []model.RssItem `xml:"channel>item"`
}

func convert(baseURL string, notes []*model.Note) []model.RssItem {
	result := make([]model.RssItem, len(notes))
	for i, n := range notes {
		result[i] = n.ToRSS(baseURL)
	}

	return result
}

func generateRss(notes []*model.Note, title, link, desc, img, rss string) ([]byte, error) {
	ba, err := xml.MarshalIndent(Rss20{
		Ns:          "http://www.w3.org/2005/Atom",
		Version:     "2.0",
		Title:       title,
		Link:        link,
		Description: desc,
		TTL:         60,
		Items:       convert(fmt.Sprintf("%s/note/", link), notes),
		Image: rssImage{
			URL:   fmt.Sprintf("%s%s", link, img),
			Title: title,
			Link:  link,
		},
		Atom: atomLink{
			Href: rss,
			Rel:  "self",
			Type: "application/rss+xml",
		},
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

	ba, err := generateRss(notes, 
		ctx.Vars["title"], 
		ctx.Vars["url"], 
		ctx.Vars["title"], 
		ctx.Vars["rss_image"],
	  ctx.Vars["rss"])
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write([]byte(xml.Header))
	w.Write(ba)

	return http.StatusOK, nil
}
