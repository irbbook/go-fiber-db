package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "localhost"
	port     = "3306"
	database = "cine_sysnc"
	username = "root"
	password = "P@ssw0rd"
)

var db *sql.DB

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	sdb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer sdb.Close()

	// Verify the connection
	if err := sdb.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	db = sdb

	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	movieGroup := app.Group("/api/movie")

	movieGroup.Get("/", getMoviesHandler)
	movieGroup.Get("/:id", getMovieByIdHandler)
	movieGroup.Post("/", createMovieHandler)
	movieGroup.Put("/:id", updateMovieHandler)
	movieGroup.Delete("/:id", deleteMovieHandler)

	log.Fatal(app.Listen(":8080"))

}

func getMoviesHandler(c *fiber.Ctx) error {
	movie, err := getMovies()

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(movie)
}

func createMovieHandler(c *fiber.Ctx) error {
	movie := new(Movie)

	if err := c.BodyParser(movie); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	err := createMovie(movie)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).SendString(err.Error())
}

func getMovieByIdHandler(c *fiber.Ctx) error {
	movieId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	movie, err := getMovie(movieId)

	if err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(movie)

}

func updateMovieHandler(c *fiber.Ctx) error {
	movieId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	movie := new(Movie)

	if err := c.BodyParser(movie); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	err = updateMovie(movieId, movie)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)

}

func deleteMovieHandler(c *fiber.Ctx) error {
	movieId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	err = deleteMovie(movieId)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)

}
