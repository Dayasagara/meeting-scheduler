package main

import (
	"log"
	"os"

	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/Dayasagara/meeting-scheduler/lib/database"
	"github.com/Dayasagara/meeting-scheduler/model"
	"github.com/Dayasagara/meeting-scheduler/receivers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	log.Println("Starting the Application at port 8080...")
	err := initGormClient()
	if err != nil {
		log.Fatalf("DB failure", err)
	}
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	g := e.Group("/calendar")
	g.GET("/ping", receivers.Get.PingHandler)
	g.POST("/signup", receivers.Post.SignUpHandler)
	g.POST("/login", receivers.Post.LoginHandler)
	g.POST("/defineAvailability", receivers.Post.DefineAvHandler)
	g.GET("/getAvailability/:date", receivers.Get.GetSlotsHandler)
	g.POST("/scheduleMeeting", receivers.Post.MeetingScheduler)
	g.GET("/getMyEvents", receivers.Get.GetEventsHandler)
	//g.POST("/syncWithGoogleCalendar", receivers.Post.SyncHandler)
	e.Start(":8080")
}

//Initialise db object(Golang's Object Relation Mapper)
func initGormClient() error {
	log.Println("Initiating DB conn")
	var config model.DBConfig
	err := godotenv.Load()
	if err != nil {
		return err
	}

	//Reads DB creds from environment variables or .env file
	config.User = os.Getenv("DBUSER")
	config.DBName = os.Getenv("DB")
	config.Password = os.Getenv("PASSWORD")
	config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")
	interfaces.DBEngine = new(database.DBRepo)
	err = interfaces.DBEngine.DBConnect(config)
	return err
}
