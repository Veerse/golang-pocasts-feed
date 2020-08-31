package app

import "time"

type Podcast struct {
	Id int
	Title,
	Description,
	Image,
	Language,
	Category,
	Author,
	Link,
	Owner string
	Episodes []Episode
}

// toFeed returns the XML of the podcast
func (p *Podcast) toFeed() string {

	return "Feed"
}

type Episode struct {
	Id          int
	Title       string
	URL         string
	Length      string
	Type        string
	Guid        int
	PubDate     time.Time
	Description string
	EpisodeURL  string
	Image       string
}
