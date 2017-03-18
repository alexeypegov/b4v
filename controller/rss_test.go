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
  
  ba, err := generateRss(notes, "title", "http://localhost/", "description")
  if err != nil {
    t.Fatal("Unable to generate RSS")
  }
  
  expected := `<rss version="2.0">
  <channel>
    <title>title</title>
    <link>http://localhost/</link>
    <description>description</description>
    <ttl>60</ttl>
    <item>
      <guid isPermaLink="true">http://localhost/локальзованный-урл</guid>
      <title>title</title>
      <category>a,b</category>
      <pubDate>Sun, 11 Nov 1979 22:23:00 +0300</pubDate>
      <description><![CDATA[Some content]]></description>
    </item>
  </channel>
</rss>`

  test.AssertEquals(expected, string(ba), t)
}