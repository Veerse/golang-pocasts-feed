package app

import "testing"

func BenchmarkPodcast_ToFeed(b *testing.B) {

}

func TestPodcast_ToFeed(t *testing.T) {
	p := Podcast{
		Id:          1,
		Title:       "La vie est belle",
		Description: "La vie est belle mais c'est une catharsis",
		Image:       "",
		Language:    "",
		Category:    "",
		AuthorName:  "",
		AuthorEmail: "",
		Link:        "",
		Owner:       "",
		Episodes:    nil,
	}

	p.ToFeed()

}
