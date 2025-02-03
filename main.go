package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WeatherInfo struct {
	ID          string  `json:id`
	Timestamp   int64   `json:timestamp`
	Temperature float32 `json:temperature`
	Humidity    float32 `json:humidity`
}

var WeatherData []WeatherInfo

func main() {
	WeatherData = []WeatherInfo{
		{ID: "1", Timestamp: 1738614276, Temperature: 19.0, Humidity: 46.0},
		{ID: "1", Timestamp: 1738614272, Temperature: 20.0, Humidity: 48.0},
	}

	router := gin.Default()
	router.GET("/weather-data", getCars)
	router.GET("/weather/:id", getCarByID)
	router.POST("/weather", createCar)
	router.DELETE("/weather", deleteCar)
	router.GET("/healthCheck", healthCheck)
	router.Run(":8080")
}

func createCar(c *gin.Context) {
	var newWeatherInfo WeatherInfo
	if err := c.BindJSON(&newWeatherInfo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to create a car"})
		return
	}
	WeatherData = append(WeatherData, newWeatherInfo)
}

func getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, WeatherData)
}

func getCarByID(c *gin.Context) {
	for _, car := range WeatherData {
		if car.ID == c.Param("id") {
			c.IndentedJSON(http.StatusOK, car)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Request car is not found"})

}

func deleteCar(c *gin.Context) {
	for index, car := range WeatherData {
		if car.ID == c.Param("id") {
			WeatherData = append(WeatherData[:index], WeatherData[index+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Car was deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car is not found"})

}

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Healthcheck is OK"})
	return
}
