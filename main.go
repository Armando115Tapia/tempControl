package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type WeatherInfo struct {
	ID          string  `json:"id"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
}

var db *sql.DB

func main() {
	// Database connection
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS weather (
			id VARCHAR(50),
			timestamp BIGINT,
			temperature FLOAT,
			humidity FLOAT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/weather-data", getWeatherData)
	router.GET("/weather/:id", getWeatherByID)
	router.POST("/weather", createWeather)
	router.DELETE("/weather/:id", deleteWeather)
	router.GET("/healthCheck", healthCheck)
	router.Run(":8080")
}

func createWeather(c *gin.Context) {
	var newWeatherInfo WeatherInfo
	if err := c.BindJSON(&newWeatherInfo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to create weather data"})
		return
	}

	_, err := db.Exec(`
		INSERT INTO weather (id, timestamp, temperature, humidity)
		VALUES ($1, $2, $3, $4)
	`, newWeatherInfo.ID, newWeatherInfo.Timestamp, newWeatherInfo.Temperature, newWeatherInfo.Humidity)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to insert data"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newWeatherInfo)
}

func getWeatherData(c *gin.Context) {
	rows, err := db.Query("SELECT id, timestamp, temperature, humidity FROM weather")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data"})
		return
	}
	defer rows.Close()

	var weatherData []WeatherInfo
	for rows.Next() {
		var w WeatherInfo
		if err := rows.Scan(&w.ID, &w.Timestamp, &w.Temperature, &w.Humidity); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan data"})
			return
		}
		weatherData = append(weatherData, w)
	}

	c.IndentedJSON(http.StatusOK, weatherData)
}

func getWeatherByID(c *gin.Context) {
	id := c.Param("id")
	var w WeatherInfo

	err := db.QueryRow("SELECT id, timestamp, temperature, humidity FROM weather WHERE id = $1", id).
		Scan(&w.ID, &w.Timestamp, &w.Temperature, &w.Humidity)

	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Weather data not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data"})
		return
	}

	c.IndentedJSON(http.StatusOK, w)
}

func deleteWeather(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM weather WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete data"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Weather data not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Weather data was deleted"})
}

func healthCheck(c *gin.Context) {
	err := db.Ping()
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "Database connection failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Healthcheck is OK"})
}
