package main

import (
	"log"
	"microservice/controllers"
	"microservice/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func main() {
	err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
		return
	}
	startEchoServer()

}

func startEchoServer() {
	e := echo.New()
	e.GET("/api/players", controllers.GetPlayers)         // done
	e.GET("/api/players/:id", controllers.GetPlayerById)  // done
	e.POST("/api/players", controllers.AddPlayer)         // done
	e.PATCH("/api/players/:id", controllers.UpdatePlayer) // done
	e.DELETE("/api/players/:id", controllers.DeletePlayer)
	e.Logger.Fatal(e.Start(":8080"))
}
