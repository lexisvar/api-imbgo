package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eefret/gomdb"
	"github.com/gin-gonic/gin"
)

func init() {
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieByID)
	router.GET("/movies/title/:title", getMoviesByTitle)
	router.GET("/movies/filter", getMoviesFilter)
	router.POST("/movies", postMovies)
	router.DELETE("/movies/:id", deleteMovieByID)
}

type cache struct {
	sync.Mutex
	movies []movie
}

// movie represents data about a record movie.
type movie struct {
	ImdbId     string   `json:"imdbid"`
	Title      string   `json:"title"`
	Released   string   `json:"released"`
	ImdbRating string   `json:"rating"`
	Genres     []string `json:"genres"`
}

// movies slice to seed record movie data.

// var movies = cache{}
var movies = cache{}

// getMovies responds with the list of all movies as JSON.
func getMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, movies.movies)
}

func getMoviesByTitle(c *gin.Context) {

	title := c.Param("title")

	// Loop over the list of movies, looking for
	// an movie whose Title value matches the parameter.
	for _, a := range movies.movies {
		if a.Title == title {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	query := &gomdb.QueryData{Title: title, SearchType: gomdb.MovieSearch}
	res, err := api.MovieByTitle(query)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "movie not found"})
		return
	}

	if res != nil {
		var tempMovie = movie{
			ImdbId:     res.ImdbID,
			Title:      res.Title,
			Released:   res.Released,
			ImdbRating: res.ImdbRating,
			Genres:     strings.Split(strings.ReplaceAll(res.Genre, " ", ""), ","), // Convert string to array
		}

		movies.Lock()
		movies.movies = append(movies.movies, tempMovie)
		movies.Unlock()

		c.IndentedJSON(http.StatusOK, tempMovie)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

// getMovies responds with the list of all movies as JSON.
func getMoviesFilter(c *gin.Context) {
	var tempMovies = []movie{}
	var start_release = c.Query("start_release")
	var end_release = c.Query("end_release")
	var rating = c.Query("rating")
	var genres = c.Query("genres")
	var option = 0 // for evaluate params

	if start_release == "" && end_release == "" && rating == "" && genres == "" {
		option = 1 // all params are empty
	} else if start_release != "" && end_release != "" {
		option = 2 // range
	} else if start_release != "" && end_release == "" {
		option = 3 // movies released same year
	}

	switch option {
	case 1:
		c.IndentedJSON(http.StatusOK, movies.movies)
		return
	case 2:
		for _, m := range movies.movies {
			if m.Released != "N/A" { // if movie have a release date

				date_realeased, err := time.Parse("2 Jan 2006", m.Released)
				if err != nil {
					c.IndentedJSON(http.StatusOK, err)
					return
				}

				int_start_release, _ := strconv.Atoi(start_release)
				int_end_release, _ := strconv.Atoi(end_release)
				if date_realeased.Year() >= int_start_release && date_realeased.Year() <= int_end_release {
					tempMovies = append(tempMovies, m)
				}
			}
		}
	case 3:
		for _, m := range movies.movies {
			if m.Released != "N/A" { // if movie have a release date

				date_realeased, err := time.Parse("2 Jan 2006", m.Released)
				if err != nil {
					c.IndentedJSON(http.StatusOK, err)
					return
				}

				int_start_release, _ := strconv.Atoi(start_release)

				if date_realeased.Year() == int_start_release {
					tempMovies = append(tempMovies, m)
				}
			}
		}
	}

	if start_release == "" && end_release == "" { // if release was not sent
		tempMovies = movies.movies
	}

	// Rating and Genres filters

	var tempMovies2 = []movie{}

	if genres != "" { // Genres
		for _, m := range tempMovies {
			if contains(m.Genres, genres) {
				tempMovies2 = append(tempMovies2, m)
			}
		}

		if len(tempMovies2) > 0 {
			tempMovies = tempMovies2
			tempMovies2 = []movie{}
		} else {
			tempMovies = []movie{}
		}
	}

	if rating != "" { // Genres
		for _, m := range tempMovies {

			if m.ImdbRating >= rating {
				tempMovies2 = append(tempMovies2, m)
			}
		}

		if len(tempMovies2) > 0 {
			tempMovies = tempMovies2
			tempMovies2 = []movie{}
		} else {
			tempMovies = []movie{}
		}
	}

	c.IndentedJSON(http.StatusOK, tempMovies)
}

// postMovies adds an movie from JSON received in the request body.
func postMovies(c *gin.Context) {
	var newMovie movie

	// Call BindJSON to bind the received JSON to
	// newMovie.
	if err := c.BindJSON(&newMovie); err != nil {
		return
	}

	// Add the new movie to the slice.
	movies.Lock()
	movies.movies = append(movies.movies, newMovie)
	movies.Unlock()
	c.IndentedJSON(http.StatusCreated, newMovie)
}

// getMovieByID locates the movie whose ID value matches the id
// parameter sent by the client, then returns that movie as a response.
func getMovieByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of movies, looking for
	// an movie whose ID value matches the parameter.
	for _, m := range movies.movies {
		if m.ImdbId == id {
			c.IndentedJSON(http.StatusOK, m)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

// deleteMovieByID delete the movie whose ID value matches the id
// parameter sent by the client, then returns that movie was deleted.
func deleteMovieByID(c *gin.Context) {
	id := c.Param("id")

	for i, m := range movies.movies {
		if m.ImdbId == id {
			movies.Lock()
			movies.movies = append(movies.movies[:i], movies.movies[i+1:]...)
			movies.Unlock()
			c.IndentedJSON(http.StatusOK, gin.H{"message": "movie deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

// functions internal
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
