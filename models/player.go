package models

import "encoding/json"

const PlayerTable = "Player"

type Player struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
	Team string `json:"team"`
}

type PlayerOutput Player

func (p *PlayerOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Age  int32  `json:"age"`
		Team string `json:"team"`
	}{
		ID:   p.ID,
		Name: p.Name,
		Age:  p.Age,
		Team: p.Team,
	})
}
