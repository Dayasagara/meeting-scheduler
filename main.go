package main

import (
	"log"

	"github.com/Dayasagara/meeting-scheduler/receivers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	log.Println("Starting the Application at port 8080...")

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.GET("/calendar/ping", receivers.Get.PingHandler)
	e.Start(":8080")
}
