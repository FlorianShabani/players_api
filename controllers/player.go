package controllers

import (
	"fmt"
	"log"
	"microservice/database"
	"microservice/models"
	"microservice/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetPlayerById(c echo.Context) error {
	id := c.Param("id")

	player, err := services.GetPlayerById(id)
	if err != nil {
		return c.JSON(404, fmt.Errorf(err.Error()))
	}

	result, err := services.UnmarshalPlayer(player)
	if err != nil {
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

func GetPlayers(c echo.Context) error {
	players, err := services.GetPlayers()
	if err != nil {
		log.Printf("Error getting players %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	result, err := services.UnmarshalPlayers(players)
	if err != nil {
		log.Printf("Error unmarshalling players %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}

func AddPlayer(c echo.Context) error {
	player := &models.Player{}

	if err := c.Bind(player); err != nil {
		log.Printf("Error binding player %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	exists, err := services.ExistsPlayer(player)

	if err != nil {
		log.Printf("Error checking player %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	if exists {
		log.Printf("Player of id %v already exists", player.ID)
		return c.JSON(409, fmt.Errorf("player of id %v already exists", player.ID))
	}

	if _, err := services.AddPlayer(player); err != nil {
		log.Printf("Error adding player %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	return c.JSON(http.StatusOK, player)
}

func UpdatePlayer(c echo.Context) error {
	id := c.Param("id")

	player := &models.Player{}

	num, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		log.Printf("Error parsing player id %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	player.ID = uint(num)

	if err := c.Bind(player); err != nil {
		log.Printf("Error binding player %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	exists, err := services.ExistsPlayer(player)

	if err != nil {
		log.Printf("Error checking player %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	if !exists {
		log.Printf("Player of id %v doesnt exist", player.ID)
		return c.JSON(404, fmt.Errorf("player of id %v already exists", player.ID))
	}

	if err := services.UpdatePlayer(player); err != nil {
		log.Printf("Error updating player %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	return c.JSON(http.StatusOK, player)
}

func DeletePlayer(c echo.Context) error {
	id := c.Param("id")

	num, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		log.Printf("Error parsing player id %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}
	//ID in uint type
	ID := uint(num)

	exists, err := database.Exists(models.PlayerTable, &ID)

	if err != nil {
		log.Printf("Error checking player %v", err.Error())
		return c.JSON(500, fmt.Errorf(err.Error()))
	}

	if !exists {
		log.Printf("Player of id %v doesnt exist", ID)
		return c.JSON(404, fmt.Errorf("player of id %v doesnt exists", ID))
	}

	if err := services.DeletePlayer(id); err != nil {
		log.Printf("Error deleting player %v", err)
		return c.JSON(209, fmt.Errorf(err.Error()))
	}
	return c.JSON(http.StatusOK, nil)
}
