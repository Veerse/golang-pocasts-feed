package app

import (
	"encoding/xml"
	"fmt"
	"github.com/eduncan911/podcast"
	"time"
)

type Podcast struct {
	Id     int
	UserId int
	Title,
	Description,
	Image,
	Language,
	Category,
	AuthorName,
	AuthorEmail,
	Link,
	Owner string
	Episodes []Episode
}

// ToFeed returns the XML of the podcast
func (p *Podcast) ToFeed() ([]byte, error) {
	feed := podcast.New(
		p.Title,
		p.Link,
		p.Description,
		nil, nil,
	)

	feed.AddAuthor(p.AuthorName, p.AuthorEmail)
	feed.AddSummary(p.Description)
	feed.AddImage(p.Image)
	feed.AddAtomLink("self")
	feed.AddCategory("Religion &amp; Spirituality", []string{"Islam"})

	feed.IOwner = &podcast.Author{
		XMLName: xml.Name{},
		Name:    p.AuthorName,
		Email:   p.AuthorEmail,
	}

	feed.IExplicit = "no"

	for _, e := range p.Episodes {
		item := podcast.Item{
			Title:       e.Title,
			Link:        e.URL,
			Description: e.Description,
			PubDate:     &e.PubDate,
			IDuration:   e.Length,
		}

		if _, err := feed.AddItem(item); err != nil {
			fmt.Printf("Adding item %+v error %s", item, err.Error())
		}
	}

	return feed.Bytes(), nil
}

type Episode struct {
	Id          int
	PodcastId   int
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

const (
	none = iota
	mosque
	imam
)

const (
	unverified = iota
	poster
	admin
)

type User struct {
	Id        int
	Type      int
	Name      string
	Street    string
	Postal    string
	City      string
	Phone     string
	Email     string
	Password  string
	Privilege int
}
