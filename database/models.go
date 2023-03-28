package database

import "time"

type UserStatsModel struct {
	ID int64

	FirstRequest     time.Time
	AmountOfRequests int64
	LastCity         string
}
