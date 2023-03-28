package executors

import (
	"errors"
	"fmt"
	"time"

	"github.com/vsabirov/tggoweatherbot/database"
)

type UserStats struct {
	ID int64

	FirstRequest     time.Time
	AmountOfRequests int64
	LastCity         string
}

func FetchUserStats(id int64) (UserStats, error) {
	model, err := database.GetUserStats(id)
	if err != nil {
		return UserStats{}, errors.New(
			fmt.Sprintf(
				"Failed to get statistics for user %d: %s",
				id, err))
	}

	return UserStats{
		ID: model.ID,

		FirstRequest:     model.FirstRequest,
		AmountOfRequests: model.AmountOfRequests,
		LastCity:         model.LastCity,
	}, nil
}

func RecordStats(id int64, city string) error {
	if !database.DoesUserStatsExist(id) {
		err := database.CreateNewUserStats(id)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to create user stats for user %d: %s", id, err))
		}

		err = database.RecordFirstUserRequest(id)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to record first user request for user %d: %s", id, err))
		}
	}

	err := database.IncreaseUserRequestAmount(id)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to increase request amount for user %d: %s", id, err))
	}

	err = database.RecordLastUserCity(id, city)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to record last requested city for user %d: %s", id, err))
	}

	return nil
}
