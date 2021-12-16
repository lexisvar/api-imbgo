package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieByID)
	router.GET("/movies/title/:title", getMoviesByTitle)
	router.POST("/movies", postMovies)
	router.DELETE("/movies/:id", deleteMovieByID)

	router.Run("localhost:9090")
}

// movie represents data about a record movie.
type movie struct {
	ImdbId   string   `json:"imdbid"`
	Title    string   `json:"title"`
	Released string   `json:"released"`
	Rated    string   `json:"artist"`
	Genres   []string `json:"genres"`
}

// movies slice to seed record movie data.

var movies = []movie{
	{ImdbId: "1", Title: "Blue Train", Released: "2005", Rated: "7.4", Genres: []string{"blah", "blah2"}},
	{ImdbId: "2", Title: "Blue Train", Released: "2005", Rated: "7.4"},
	{ImdbId: "3", Title: "Blue Train", Released: "2005", Rated: "7.4"},
}

// getMovies responds with the list of all movies as JSON.
func getMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, movies)
}

func getMoviesByTitle(c *gin.Context) {
	// api := gomdb.Init("9b41e7cc")
	// query := &gomdb.QueryData{Title: "Macbeth", SearchType: gomdb.MovieSearch}
	// res, err := api.MovieByTitle(query)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusOK, err)
	// 	return
	// }
	// c.IndentedJSON(http.StatusOK, res)
	title := c.Param("title")

	// Loop over the list of movies, looking for
	// an movie whose ID value matches the parameter.
	for _, a := range movies {
		if a.Title == title {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
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
	movies = append(movies, newMovie)
	c.IndentedJSON(http.StatusCreated, newMovie)
}

// getMovieByID locates the movie whose ID value matches the id
// parameter sent by the client, then returns that movie as a response.
func getMovieByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of movies, looking for
	// an movie whose ID value matches the parameter.
	for _, a := range movies {
		if a.ImdbId == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

// deleteMovieByID delete the movie whose ID value matches the id
// parameter sent by the client, then returns that movie was deleted.
func deleteMovieByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range movies {
		if a.ImdbId == id {
			movies = append(movies[:i], movies[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "movie deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}
