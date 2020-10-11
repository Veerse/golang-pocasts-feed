package app

import (
	"encoding/xml"
	"github.com/eduncan911/podcast"
	"time"
)

// See documentation https://itunespartner.apple.com/podcasts/articles/podcast-requirements-3058

type Podcast struct {
	Id          int       `json:"id"`
	UserId      int       `form:"user_id" json:"user_id" binding:"-"`
	Title       string    `form:"title" json:"title" binding:"required"`
	Description string    `form:"description" json:"description" binding:"required"`
	Image       string    `form:"image" json:"image" binding:"required"`
	Language    string    `form:"language" json:"language" binding:"required"`
	Category    string    `form:"category" json:"category" binding:"required"`
	AuthorName  string    `form:"author_name" json:"author_name" binding:"required"`
	AuthorEmail string    `form:"author_email" json:"author_email" binding:"required"`
	Link        string    `form:"link" json:"link" binding:"required"`
	Owner       string    `form:"owner" json:"owner" binding:"required"`
	Episodes    []Episode `json:"episodes" binding:"-"`
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
			return feed.Bytes(), err
		}
	}

	return feed.Bytes(), nil
}

type Episode struct {
	Id          int       `json:"id" binding:"-"`
	PodcastId   int       `form:"podcast_id" json:"podcast_id" binding:"-"`
	Title       string    `form:"title" json:"title" binding:"required"`
	URL         string    `form:"url" json:"url" binding:"required"`
	Length      string    `form:"length" json:"length" binding:"required"`
	Type        string    `form:"type" json:"type" binding:"required"`
	Guid        int       `form:"guid" json:"guid" binding:"required"`
	PubDate     time.Time `form:"id" json:"pub_date" binding:"required"`
	Description string    `form:"description" json:"description" binding:"required"`
	EpisodeURL  string    `form:"episode_url" json:"episode_url" binding:"required"`
	Image       string    `form:"image" json:"image" binding:"required"`
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
