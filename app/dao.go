package app

import (
	"database/sql"
)

func GetAllPodcastsDao(db *sql.DB) ([]Podcast, error) {
	var podcasts []Podcast

	rows, err := db.Query("select id, title, description, image, language, category, author_name, author_email, link, owner from podcasts")

	if err != nil {
		return podcasts, err
	}

	for rows.Next() {
		var p Podcast

		if err := rows.Scan(&p.Id, &p.Title, &p.Description, &p.Image, &p.Language, &p.Category, &p.AuthorName, &p.AuthorEmail, &p.Link, &p.Owner); err != nil {
			return podcasts, err
		} else {
			if rows, err := db.Query("select id, title, url, length, type, guid, pub_date, description, episode_url, image from episodes where podcast_id = $1", p.Id); err != nil {
				return podcasts, err
			} else {
				for rows.Next() {
					var e Episode
					if err2 := rows.Scan(&e.Id, &e.Title, &e.URL, &e.Length, &e.Type, &e.Guid, &e.PubDate, &e.Description, &e.EpisodeURL, &e.Image); err2 != nil {
						return podcasts, err2
					} else {
						p.Episodes = append(p.Episodes, e)
					}
				}
			}
			podcasts = append(podcasts, p)
		}
	}

	return podcasts, nil
}

func GetPodcastByIdDao(id string, db *sql.DB) (Podcast, error) {
	var podcast Podcast

	if err := db.QueryRow("select id, title, description, image, language, category, author_name, author_email, link, owner from podcasts where id = $1", id).
		Scan(&podcast.Id, &podcast.Title, &podcast.Description, &podcast.Image, &podcast.Language,
			&podcast.Category, &podcast.AuthorName, &podcast.AuthorEmail, &podcast.Link, &podcast.Owner); err != nil {
		return podcast, err
	}

	if rows, err := db.Query("select id, title, url, length, type, guid, pub_date, description, episode_url, image from episodes where podcast_id = $1", id); err != nil {
		return podcast, err
	} else {
		for rows.Next() {
			var e Episode
			if err2 := rows.Scan(&e.Id, &e.Title, &e.URL, &e.Length, &e.Type, &e.Guid, &e.PubDate, &e.Description, &e.EpisodeURL, &e.Image); err2 != nil {
				return podcast, err2
			} else {
				podcast.Episodes = append(podcast.Episodes, e)
			}
		}
	}

	return podcast, nil
}
