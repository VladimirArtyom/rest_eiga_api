package data

import (
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

type Movie struct {
	ID      int64    `json:"id"` //Unique identifier for the movei
	Title   string   `json:"title"`
	Year    int32    `json:"year,omitempty"`
	Runtime Runtime  `json:"runtime,omitempty"` //in mins
	Genres  []string `json:"genres,omitempty"`
	Version int32    `json:"version"`

	CreatedAt time.Time `json:"-"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {

	//validate the input

	//Please separate each error, don't mix it
	// Add errors if exists, only invalid check is added.
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must lower than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year > 1888, "year", "must be greater than 1888")
	v.Check(movie.Year < int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) > 0, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) < 6, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

}
