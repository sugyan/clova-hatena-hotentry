package hatena

import (
	"encoding/xml"
	"errors"
	"net/http"
	"time"
)

const urlBase = "http://b.hatena.ne.jp"

// Category type
type Category string

// Constant values
const (
	CategoryAll           Category = "総合"
	CategoryGeneral       Category = "一般"
	CategorySocial        Category = "世の中"
	CategoryEconomics     Category = "政治と経済"
	CategoryLife          Category = "暮らし"
	CategoryKnowledge     Category = "学び"
	CategoryIt            Category = "テクノロジー"
	CategoryEntertainment Category = "エンタメ"
	CategoryGame          Category = "アニメとゲーム"
	CategoryFun           Category = "おもしろ"
)

var categoryMap = map[Category]string{
	CategoryAll:           "/hotentry.rss",
	CategoryGeneral:       "/hotentry/general.rss",
	CategorySocial:        "/hotentry/social.rss",
	CategoryEconomics:     "/hotentry/economics.rss",
	CategoryLife:          "/hotentry/life.rss",
	CategoryKnowledge:     "/hotentry/knowledge.rss",
	CategoryIt:            "/hotentry/it.rss",
	CategoryEntertainment: "/hotentry/entertainment.rss",
	CategoryGame:          "/hotentry/game.rss",
	CategoryFun:           "/hotentry/fun.rss",
}

// Entry type
type Entry struct {
	Title         string    `xml:"title"`
	BookmarkCount int       `xml:"bookmarkcount"`
	Date          time.Time `xml:"date"`
	Subjects      []string  `xml:"subject"`
}

// Client type
type Client struct {
	httpClient *http.Client
}

// ClientOption type
type ClientOption func(*Client)

// NewClient function
func NewClient(options ...ClientOption) *Client {
	c := &Client{
		httpClient: http.DefaultClient,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// WithHTTPClient function
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = c
	}
}

// Fetch method
func (c *Client) Fetch(category Category) ([]*Entry, error) {
	path, ok := categoryMap[category]
	if !ok {
		return nil, errors.New("Category not found")
	}

	res, err := c.httpClient.Get(urlBase + path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Title string   `xml:"channel>title"`
		Item  []*Entry `xml:"item"`
	}
	if err := xml.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Item, nil
}
