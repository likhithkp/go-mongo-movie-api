package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(c *gin.Context) {
	c.JSON(http.StatusOK, movies)
}

func getMovie(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Please enter the ID"})
		return
	}

	for _, movie := range movies {
		if movie.ID == string(id) {
			c.JSON(http.StatusOK, movie)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
}

func addMovie(c *gin.Context) {
	var newMovie Movie

	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newMovie.ID = strconv.Itoa(len(movies) + 1)
	movies = append(movies, newMovie)
	c.JSON(http.StatusCreated, gin.H{"movie": newMovie})
}

func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var updatedMovie Movie

	if err := c.ShouldBindJSON(&updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for index, movie := range movies {
		if movie.ID == id {
			movies[index].Title = updatedMovie.Title
			movies[index].Director = updatedMovie.Director

			c.JSON(http.StatusOK, gin.H{"message": "Movie updated", "data": movies[index]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
}

func deleteMovie(c *gin.Context) {
	id := c.Param("id")
	for index, movie := range movies {
		if movie.ID == string(id) {
			movies = append(movies[:index], movies[index+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Movie" + movie.Title + "deleted",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to delet"})
	}
}

func main() {
	router := gin.Default()

	movies = append(movies, Movie{ID: "1", Title: "Rio", Director: &Director{Firstname: "Shiba", Lastname: "Kumar"}})
	movies = append(movies, Movie{ID: "2", Title: "Rio 2", Director: &Director{Firstname: "Shiba", Lastname: "Kumar"}})

	router.GET("/allMovies", getMovies)
	router.GET("/getMovie/:id", getMovie)
	router.POST("/addMovie", addMovie)
	router.PUT("/updateMovie/:id", updateMovie)
	router.DELETE("/deleteMovie/:id", deleteMovie)

	http.ListenAndServe(":3000", router)
}
