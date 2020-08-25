package app

import (
	"database/sql"
)

func GetAllPodcastsDao(db *sql.DB) ([]Podcast, error) {
	var podcasts []Podcast

	return podcasts, nil
}

func GetPodcastByIdDao(id string, db *sql.DB) (Podcast, error) {
	var podcast Podcast

	if err := db.QueryRow("select id, title, description, image, language, category, author, link, owner from podcast where id = $1", id).
		Scan(&podcast.Id, &podcast.Title, &podcast.Description, &podcast.Image, &podcast.Language,
			&podcast.Category, &podcast.Author, &podcast.Link, &podcast.Owner); err != nil {
		return podcast, err
	}

	return podcast, nil
}