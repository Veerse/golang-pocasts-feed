package app

import (
	"database/sql"
)

func GetAllPodcastsDao(db *sql.DB) ([]Podcast, error) {
	var podcasts []Podcast

	rows, err := db.Query("select id, user_id, title, description, image, language, category, author_name, author_email, link, owner from podcasts")

	if err != nil {
		return podcasts, err
	}

	for rows.Next() {
		var p Podcast

		if err := rows.Scan(&p.Id, &p.UserId, &p.Title, &p.Description, &p.Image, &p.Language, &p.Category, &p.AuthorName, &p.AuthorEmail, &p.Link, &p.Owner); err != nil {
			return podcasts, err
		} else {
			if rows, err := db.Query("select id, podcast_id, title, url, length, type, guid, pub_date, description, episode_url, image from episodes where podcast_id = $1", p.Id); err != nil {
				return podcasts, err
			} else {
				for rows.Next() {
					var e Episode
					if err2 := rows.Scan(&e.Id, &e.PodcastId, &e.Title, &e.URL, &e.Length, &e.Type, &e.Guid, &e.PubDate, &e.Description, &e.EpisodeURL, &e.Image); err2 != nil {
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
	var p Podcast

	if err := db.QueryRow("select id, user_id, title, description, image, language, category, author_name, author_email, link, owner from podcasts where id = $1", id).
		Scan(&p.Id, &p.UserId, &p.Title, &p.Description, &p.Image, &p.Language,
			&p.Category, &p.AuthorName, &p.AuthorEmail, &p.Link, &p.Owner); err != nil {
		return p, err
	}

	if rows, err := db.Query("select id, podcast_id, title, url, length, type, guid, pub_date, description, episode_url, image from episodes where podcast_id = $1", id); err != nil {
		return p, err
	} else {
		for rows.Next() {
			var e Episode
			if err2 := rows.Scan(&e.Id, &e.PodcastId, &e.Title, &e.URL, &e.Length, &e.Type, &e.Guid, &e.PubDate, &e.Description, &e.EpisodeURL, &e.Image); err2 != nil {
				return p, err2
			} else {
				p.Episodes = append(p.Episodes, e)
			}
		}
	}

	return p, nil
}

func GetUserByIdDao(id string, db *sql.DB) (User, error) {
	return User{}, nil
}

func GetUserByEmailAndPassword(email, password string, db *sql.DB) (User, error) {
	var u User

	if err := db.QueryRow("select id, type, name, street, postal, city, phone, email, password, privilege from Users where email = $1 and password = $2", email, password).
		Scan(&u.Id, &u.Type, &u.Name, &u.Street, &u.Postal, &u.City, &u.Phone, &u.Email, &u.Password, &u.Privilege); err != nil {
		return u, err
	}
	return u, nil
}
