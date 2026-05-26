package rss

import (
	"crypto/md5"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stefannegele/website/internal/db"
)

const innoqFeedURL = "https://www.innoq.com/en/written.atom"
const authorEmail = "stefan.negele@innoq.com"

func FetchInnoqPosts() ([]db.Post, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(innoqFeedURL)
	if err != nil {
		return nil, fmt.Errorf("fetching INNOQ feed: %w", err)
	}

	var posts []db.Post
	for _, item := range feed.Items {
		if !isAuthor(item) {
			continue
		}

		published := time.Now()
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			published = *item.UpdatedParsed
		}

		slug := slugify(item.Title) + "-" + hashShort(item.Link)
		summary := stripHTML(item.Description)
		if len(summary) > 300 {
			summary = summary[:300] + "…"
		}

		posts = append(posts, db.Post{
			Title:       item.Title,
			Slug:        slug,
			Summary:     summary,
			Content:     item.Content,
			URL:         item.Link,
			Source:      "innoq",
			PublishedAt: published,
		})
	}
	return posts, nil
}

func isAuthor(item *gofeed.Item) bool {
	for _, a := range item.Authors {
		if strings.EqualFold(a.Email, authorEmail) ||
			strings.EqualFold(a.Name, "Stefan Negele") {
			return true
		}
	}
	return false
}

func slugify(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func hashShort(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[:6]
}

func stripHTML(s string) string {
	s = html.UnescapeString(s)
	re := regexp.MustCompile(`<[^>]+>`)
	return strings.TrimSpace(re.ReplaceAllString(s, ""))
}
