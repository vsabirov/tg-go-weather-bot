package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Credentials struct {
	Host string
	Port int

	User     string
	Password string

	Database string
}

var database *sql.DB

func Connect(credentials Credentials) error {
	destination := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		credentials.User, credentials.Password, credentials.Host, credentials.Port, credentials.Database)

	connection, err := sql.Open("postgres", destination)
	if err != nil {
		return err
	}

	database = connection

	database.SetConnMaxLifetime(time.Minute * 3)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)

	return nil
}

func Disconnect() error {
	if database == nil {
		return errors.New("Database was not opened.")
	}

	database.Close()

	return nil
}

func CreateNewUserStats(id int64) error {
	insert := `INSERT INTO UserStats (id) VALUES ($1);`
	_, err := database.Exec(insert, id)

	return err
}

func IncreaseUserRequestAmount(id int64) error {
	update := `UPDATE UserStats SET amountOfRequests = amountOfRequests + 1 WHERE id = $1;`
	_, err := database.Exec(update, id)

	return err
}

func RecordFirstUserRequest(id int64) error {
	update := `UPDATE UserStats SET firstRequest = NOW() WHERE id = $1;`
	_, err := database.Exec(update, id)

	return err
}

func RecordLastUserCity(id int64, city string) error {
	update := `UPDATE UserStats SET lastCity = $1 WHERE id = $2;`
	_, err := database.Exec(update, city, id)

	return err
}

func GetUserStats(id int64) (UserStatsModel, error) {
	selection := `SELECT * FROM UserStats WHERE id = $1;`
	rows, err := database.Query(selection, id)
	if err != nil {
		return UserStatsModel{}, err
	}

	defer rows.Close()

	result := UserStatsModel{}
	if rows.Next() {
		err = rows.Scan(&result.ID, &result.FirstRequest, &result.AmountOfRequests, &result.LastCity)
		if err != nil {
			return UserStatsModel{}, err
		}

		return result, nil
	}

	return UserStatsModel{}, errors.New("No row found.")
}

func DoesUserStatsExist(id int64) bool {
	_, err := GetUserStats(id)

	return err == nil
}
