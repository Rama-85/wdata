package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WeatherData struct {
	//gorm.Model
	//WisId              uint    `json:"wis_id" gorm:"primary_key"`
	DeviceNum  string `json:"device_num" binding:"required" `
	DeviceName string `json:"device_name" binding:"required"`
	DateTime   string `json:"date_time" binding:"required"`
	//StartDate          string  `json:"start_date" binding:"required"`
	//EndDate            string  `json:"end_date" binding:"required"`
	StationName        string  `json:"station_name" binding:"required"`
	AirTemperature     float64 `json:"air_temperature" binding:"required"`
	BatteryValues      float64 `json:"battery_values" binding:"required"`
	RelativeHumidity   float64 `json:"relative_humidity" binding:"required"`
	RoadTemperature    float64 `json:"road_temperature" binding:"required"`
	Visibility         float64 `json:"visibility" binding:"required"`
	WindDirection      float64 `json:"wind_direction" binding:"required"`
	WindSpeed          int     `json:"wind_speed" binding:"required"`
	WindSpeedKmh       float64 `json:"wind_speed_kmh" binding:"required"`
	Rain               float64 `json:"rain" binding:"required"`
	WindCardinal       string  `json:"wind_cardinal" binding:"required"`
	AtmospherePressure float64 `json:"atmosphere_pressure" binding:"required"`
	IrstPav            float64 `json:"irst_pav" binding:"required"`
	Dateadded          string  `json:"date_added" binding:"required"`
}

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/itms"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&WeatherData{})

	DB = db

}

func main() {
	ConnectDatabase()

	router := gin.Default()
	//router.GET("/weather_data", GetAllData)

	router.GET("/weather_data", func(c *gin.Context) {

		device_num := c.DefaultQuery("device_num", "4242")
		//date_time :=c.query("date_time")
		start_date_query := c.DefaultQuery("start_date", "2023-08-17")
		end_date_query := c.DefaultQuery("end_date", "2023-08-17")

		start_date_query_parse, err := time.Parse("2006-01-02 15:04:05", start_date_query)
		if err != nil {
			fmt.Println("while parse start date : ", err)
			return
		}
		end_date_query_parse, err := time.Parse("2006-01-02 15:04:05", end_date_query)
		if err != nil {
			fmt.Println("while parse end date : ", err)
			return
		}

		var weather_data []WeatherData
		// Order("date_time asc")
		err = DB.Find(&weather_data).Error
		if err != nil {
			fmt.Println("err while fetch : ", err)
			return
		}
		var weather_data_result []WeatherData
		//var device_num, date_time string
		for _, data := range weather_data {

			date_date_time, err := time.Parse("2006-01-02 15:04:05", data.DateTime)
			if err != nil {
				fmt.Println("while parse start date : ", err)
				return
			}

			// && start_dateStr.Format("2023-08-17 14:31:00") == start_date && end_dateStr.Format("2023-08-17 14:31:16") == end_date
			if (data.DeviceNum) == device_num &&
				(date_date_time.Equal(start_date_query_parse) || date_date_time.After(start_date_query_parse)) &&
				(date_date_time.Equal(end_date_query_parse) || date_date_time.Before(end_date_query_parse)) {
				weather_data_result = append(weather_data_result, data)
			}
		}
		if device_num == "" {
			weather_data_result = weather_data
		}

		c.JSON(http.StatusOK, gin.H{"message": "This is for weather_data", "data": weather_data_result})
	})
	router.Run(":8080")
}
