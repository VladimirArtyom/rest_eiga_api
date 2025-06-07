package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"` //Unique identifier for the movei
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` //in mins
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"-"`
}

type MovieModel struct {
	DB *sql.DB
}

func (m *MovieModel) Insert(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var args []interface{}
	args = append(args, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)) // Make sure that each datatype has been supported by the database to read.

	// Save the returning variables to existing movie.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m *MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector(title) @@ plainto_tsquery('simple', $1) OR $1 = '')  
		AND (genres @> $2 OR $2 = ARRAY[]::TEXT[] )
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
	`, filters.sortColumn(),
		filters.sortDirection()) // Using no stemming approach for tsquery

	// creer une context avec 3-seconds timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		title,
		pq.Array(genres),
		filters.limit(),
		filters.offset(),
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	movies := []*Movie{}
	var totalRecords int = 0
	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	metadata := calculateMetadata(totalRecords,
		filters.Page,
		filters.PageSize)

	if rows.Err() != nil {
		return nil, Metadata{}, err
	}

	return movies, metadata, nil
}

func (m *MovieModel) Get(id int64) (*Movie, error) {

	/*
		In case tu peux utiliser ça
		query := `
			SELECT pg_sleep(10), id, created_at, title, year, runtime, genres, version
			FROM movies
			WHERE id = $1
		`
		var movie Movie = Movie{}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		// La méthode cancel est appelée avant le retour de la fonction.
		defer cancel()

		err := m.DB.QueryRowContext(ctx, query, id).Scan(
			&[]byte{},
			&movie.ID, &movie.CreatedAt,
			&movie.Title, &movie.Year, &movie.Runtime,
			pq.Array(&movie.Genres), &movie.Version)
	*/
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		from movies
		WHERE id=$1
	`
	var movie Movie = Movie{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID, &movie.CreatedAt,
		&movie.Title, &movie.Year, &movie.Runtime,
		pq.Array(&movie.Genres), &movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	var query string = `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *MovieModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM movies
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sqlResult, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return ErrRecordNotFound
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return ErrRecordNotFound
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
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
