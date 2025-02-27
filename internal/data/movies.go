package data

import "time"

type Movie struct {
	ID int64 `json:"id"` //Unique identifier for the movei
	Title string `json:"title"`
	Year int32 `json:"year,omitempty"`
	Runtime Runtime `json:"runtime,omitempty"` //in mins
	Genres []string `json:"genres,omitempty"`
	Version int32 `json:"version"`

	CreatedAt time.Time `json:"-"`
}


