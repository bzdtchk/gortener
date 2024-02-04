package database

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type ShortenedLink struct {
	Id             int       `json:"id"`
	Slug           string    `json:"slug"`
	DestinationUrl string    `json:"destination_url"`
	CreatedAt      time.Time `json:"created_at"`
}

func GetBySlug(slug string) (*ShortenedLink, error) {
	db, _ = GetDB()

	row := db.QueryRow("SELECT * FROM links WHERE slug = ?", slug)

	shortenedLink := ShortenedLink{}
	var err error
	if err = row.Scan(&shortenedLink.Id, &shortenedLink.Slug, &shortenedLink.DestinationUrl, &shortenedLink.CreatedAt); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return nil, errors.New("Link not found")
	}
	return &shortenedLink, nil
}

func InsertLink(link *ShortenedLink) (*ShortenedLink, error) {
	db, _ = GetDB()

	res, err := db.Exec("INSERT INTO links (slug, destination_url) VALUES (?,?)", link.Slug, link.DestinationUrl)
	if err != nil {
		return &ShortenedLink{}, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return &ShortenedLink{}, err
	}
	link.Id = int(id)
	return link, nil
}
