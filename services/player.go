package services

import (
	"encoding/json"
	"errors"
	"log"
	"microservice/database"
	"microservice/models"
	"strconv"
)

func GetPlayerById(id string) (*models.Player, error) {
	player := &models.Player{}

	sqlID, err := strconv.Atoi(id)

	if err != nil {
		log.Printf("Error during conversion  %s", err)
		return nil, err
	}
	row, err := database.QueryOneById(models.PlayerTable, &sqlID)
	if err != nil {
		log.Printf("Error getting player  %s", err)
		return nil, errors.New("could not find player")
	}

	err = row.Scan(&player.ID, &player.Name, &player.Age, &player.Team)

	if err != nil {
		log.Printf("Error scanning row %s", err)
		return nil, errors.New("could not scan player")
	}

	return player, nil
}

func UnmarshalPlayer(p *models.Player) (*map[string]interface{}, error) {
	output := models.PlayerOutput(*p)

	js, err := json.Marshal(&output)
	if err != nil {
		log.Printf("func UnmarshalPlayer: error marshalling  %s", err)
		return nil, errors.New("could not umarshal player")
	}

	result := make(map[string]interface{})

	if err := json.Unmarshal(js, &result); err != nil {
		log.Printf("func UnmarshalPlayer: error in creating player %s", err)
		return nil, err
	}

	return &result, nil
}

func GetPlayers() (*[]models.Player, error) {
	rows, err := database.Query(models.PlayerTable)

	if err != nil {
		log.Printf("func GetPlayers: error %s", err)
		return nil, errors.New("error getting players")
	}

	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		err := rows.Scan(&player.ID, &player.Name, &player.Age, &player.Team)
		if err != nil {
			log.Printf("func_QueryDatabase: Error when scanning database %s", err)
			return nil, errors.New("QueryDatabase: Error when scanning")
		}

		players = append(players, player)
	}
	if err := rows.Err(); err != nil { //Qa po bon qekjo?
		log.Printf("func_QueryDatabase: Error in rows.Err %s", err)
		return nil, errors.New("QueryDatabase: Error when scanning")
	}

	return &players, nil
}

func UnmarshalPlayers(players *[]models.Player) (*[]map[string]interface{}, error) {
	result := [](map[string]interface{}){}
	for _, p := range *players {
		v, err := UnmarshalPlayer(&p)
		if err != nil {
			log.Printf("func_UnmarshalPlayers Error unmarshalling player %s", err)
			return nil, err
		}
		result = append(result, *v)
	}
	return &result, nil
}

func AddPlayer(p *models.Player) (*models.Player, error) {
	p, err := database.CreatePlayer(p)
	if err != nil {
		log.Printf("func services.AddPlayer: error adding player to database %s", err)
		return nil, err
	}
	return p, nil
}

func ExistsPlayer(p *models.Player) (bool, error) {
	return database.Exists(models.PlayerTable, &p.ID)
}

func UpdatePlayer(p *models.Player) error {
	err := database.UpdatePlayer(p)
	if err != nil {
		log.Printf("func services.UpdatePlayer: error updating player in database %s", err)
		return err
	}
	return nil
}

func DeletePlayer(id string) error {
	err := database.DeleteRecordById(models.PlayerTable, id)
	if err != nil {
		log.Printf("func services.DeletePlayer: error deleting player in database %s", err)
		return err
	}
	return nil
}
