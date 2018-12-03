package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type Datastore interface {
	GetTrips(medallion []string, date string) (map[string]string, error)
}

type DB struct {
	*sql.DB
}

const (
	table = "cab_trip_data"
)

func NewDB(dataSource string) (*DB, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) GetTrips(medallion []string, date string) (map[string]string, error) {
	datStr := date + "%"
	result := make(map[string]string)
	var count int
	for i := 0; i < len(medallion); i++ {
		count = 0
		err := db.QueryRow("select count(*) from cab_trip_data where medallion IN(?) and pickup_datetime LIKE(?)", medallion[i], datStr).Scan(&count)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error: Query failed for medallion:%v - :%v", medallion[i], err))
		} else if count == 0 {
			results, err := db.Query("select * from cab_trip_data where medallion IN(?) and pickup_datetime LIKE(?)", medallion[i], datStr)
			if err != nil || results.Err() == nil {
				result[medallion[i]] = "Data not found."
				continue
			}
		}
		result[medallion[i]] = strconv.Itoa(count)
	}

	return result, nil
}
