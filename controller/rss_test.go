package controller

import (
  "testing"
  "time"
  
  "github.com/alexeypegov/b4v/model"
  "github.com/alexeypegov/b4v/test"
)

func TestRss(t *testing.T) {
  ts, _ := time.Parse(time.RFC822, "11 Nov 79 22:23 MSK")

  note := model.Note{
    UUID:      "локальзованный-урл",
    Title:     "title",
    Content:   "Some content",
    CreatedAt: ts,
    Tags:      []string{"a", "b"},
  }

  notes := make([]*model.Note, 1)
  notes[0] = &note
  
  ba, err := generateRss(notes, "title", "http://localhost", "description", "/img/favicon.png", "/rss")
  if err != nil {
    t.Fatal("Unable to generate RSS")
  }
  
  expected := `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>title</title>
    <link>http://localhost</link>
    <description>description</description>
    <ttl>60</ttl>
    <image>
      <url>http://localhost/img/favicon.png</url>
      <title>title</title>
      <link>http://localhost</link>
    </image>
    <atom:link href="/rss" rel="self" type="application/rss+xml"></atom:link>
    <item>
      <guid isPermaLink="true">http://localhost/note/локальзованный-урл</guid>
      <title>title</title>
      <category>a,b</category>
      <pubDate>Sun, 11 Nov 1979 22:23:00 +0300</pubDate>
      <description><![CDATA[Some content]]></description>
    </item>
  </channel>
</rss>`

  test.AssertEquals(expected, string(ba), t)
}