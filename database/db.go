package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"microservice/models"
)

var (
	config *Config
)

const DBName = "Players_API"

type Config struct {
	MySQL struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"mysql"`
}

func loadConfig() {
	//Get username and password of database connection from config file
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error reading config.json")
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error unmarshalling config.json")
	}

}

func InitDatabase() error {
	loadConfig()
	db, err := OpenDatabase()
	if err != nil {
		log.Fatal("Error when initializing database", err)
		return errors.New("Error opening database on init")
	}
	defer db.Close()

	createPlayersTable := "CREATE TABLE IF NOT EXISTS " + models.PlayerTable + "(     id INT PRIMARY KEY,     name VARCHAR(255) NOT NULL,     age INT NOT NULL,     team VARCHAR(255) NOT NULL );"
	_, err = db.Exec(createPlayersTable)

	if err != nil {
		log.Fatal("Error when executing create players table script", err)
		return errors.New("Error creating players table")
	}

	return nil
}

func OpenDatabase() (*sql.DB, error) {
	//Open database connection
	database, err := sql.Open("mysql", config.MySQL.Username+":"+config.MySQL.Password+"@tcp(127.0.0.1:3306)/"+DBName)

	if err != nil {
		log.Fatal("func_OpenDatabase: Error when opening database", err)
		return nil, errors.New("sadfasdf")
	}

	return database, nil
}

func Query(table string) (*sql.Rows, error) {
	db, err := OpenDatabase()
	if err != nil {
		log.Fatal("func_Query Error when opening database on query", err)
		return nil, errors.New("Error when openin database")
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		log.Printf("func_Query Error when querying database %s", err)
		return nil, errors.New("Query: Error when querying")
	}

	return rows, nil
}

// identifier as : ID
func QueryOneById(table string, id *int) (*sql.Row, error) {
	db, err := OpenDatabase()
	if err != nil {
		log.Printf("func_QueryOne: Error when opening database %s", err)
		return nil, errors.New("Query: Error opening database")
	}

	defer db.Close()

	//row := db.QueryRow("SELECT * FROM " + table + " WHERE " + *identifier)
	row := db.QueryRow("SELECT * FROM "+table+" WHERE id = ?", *id)
	if err := row.Scan(); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no player found with id %v", *id)
			return nil, errors.New("Query: No row with ID")
		} else {
			log.Printf("func_QueryOne: Error when querying database %s", err)
			return nil, errors.New("Query: Error when querying")
		}
	}

	return row, nil
}

func CreatePlayer(p *models.Player) error {
	db, err := OpenDatabase()
	if err != nil {
		log.Printf("func_CreatePlayer: Error when opening database %s", err)
		return errors.New("error when opening database")
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO "+models.PlayerTable+" (id, name, age, team) VALUES (?, ?,?,?)", p.ID, p.Name, p.Age, p.Team)
	if err != nil {
		log.Printf("func_CreatePlayer: Error when creating player %s", err)
		return errors.New("error when creating player")
	}

	return nil
}

func Exists(table string, id *uint) (bool, error) {
	db, err := OpenDatabase()
	if err != nil {
		log.Printf("func_CreatePlayer: Error when opening database %s", err)
		return false, errors.New("error when opening database")
	}

	defer db.Close()

	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM "+table+" WHERE id = ?)", id)

	var exists bool = true

	if err := row.Scan(&exists); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("func_QueryOne: Error when querying database %s", err)
			return false, errors.New("Query: Error when querying")
		}
	}
	return exists, nil
}

func UpdatePlayer(p *models.Player) error {

	/*
			UPDATE employees
		SET
		    email = 'mary.patterson@classicmodelcars.com'
		WHERE
		    employeeNumber = 1056;
	*/

	db, err := OpenDatabase()
	if err != nil {
		log.Printf("func_UpdatePlayer: Error when opening database %s", err)
		return errors.New("error when opening database")
	}

	defer db.Close()

	if p.Name != "" {
		_, err = db.Exec("UPDATE "+models.PlayerTable+" SET name =? WHERE id =?", p.Name, p.ID)
		if err != nil {
			log.Printf("func_UpdatePlayer: Error when updating player name %s", err)
			return errors.New("error when updating player")
		}
	}

	if p.Age != 0 {
		_, err = db.Exec("UPDATE "+models.PlayerTable+" SET age =? WHERE id =?", p.Age, p.ID)
		if err != nil {
			log.Printf("func_UpdatePlayer: Error when updating player age %s", err)
			return errors.New("error when updating player")
		}
	}

	if p.Team != "" {
		_, err = db.Exec("UPDATE "+models.PlayerTable+" SET team =? WHERE id =?", p.Team, p.ID)
		if err != nil {
			log.Printf("func_UpdatePlayer: Error when updating player team %s", err)
			return errors.New("error when updating player")
		}
	}

	return nil
}

func DeleteRecordById(table string, id string) error {
	//DELETE FROM player WHERE id = 123;

	db, err := OpenDatabase()
	if err != nil {
		log.Printf("DeletePlayer: Error when opening database %s", err)
		return errors.New("error when opening database")
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM " + table + " WHERE id = " + id)
	if err != nil {
		log.Printf("DeletePlayer: Error when deleting player %s", err)
		return errors.New("error when creating player")
	}

	return nil
}
