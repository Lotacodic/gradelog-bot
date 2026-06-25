package storage

import (
	"encoding/json"
	"os"

	"gradelog-sys/models"
)

const DBPath = "database.json"

func LoadDB() (models.Database, error) {
	data, err := os.ReadFile(DBPath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.Database{Students: []models.Student{}}, nil
		}
		return models.Database{}, err
	}

	var db models.Database
	err = json.Unmarshal(data, &db)
	if err != nil {
		return models.Database{}, err
	}

	return db, nil
}

func SaveDB(db models.Database) error {
	data, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(DBPath, data, 0o644)
}
