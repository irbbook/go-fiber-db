package main

func createMovie(movie *Movie) error {
	_, err := db.Exec(
		"INSERT INTO movies (title, director, year) VALUES (?, ?, ?);",
		movie.Title, movie.Director, movie.Year,
	)
	return err
}

func getMovie(id int) (Movie, error) {
	var m Movie
	row := db.QueryRow("SELECT id,title, director, year FROM movies WHERE id = ?", id)

	err := row.Scan(&m.ID, &m.Title, &m.Director, &m.Year)

	if err != nil {
		return Movie{}, err
	}

	return m, nil
}

func getMovies() ([]Movie, error) {
	rows, err := db.Query("SELECT id,title, director, year FROM movies;")

	if err != nil {
		return []Movie{}, err
	}

	var movies []Movie

	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Director, &m.Year)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func updateMovie(id int, movie *Movie) error {
	_, err := db.Exec(
		"UPDATE movies SET title = ?, director = ? , year = ? WHERE id = ?;",
		movie.Title,
		movie.Director,
		movie.Year,
		id,
	)
	return err
}

func deleteMovie(id int) error {
	_, err := db.Exec("DELETE FROM movies WHERE id = ?;", id)
	return err
}
