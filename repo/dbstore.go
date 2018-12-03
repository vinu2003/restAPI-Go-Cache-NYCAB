package repo

import (
	"database/sql"
	"log"
)

type Datastore interface {
	GetTrips(medallion []string, date string) (map[string]int, error)
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

func (db *DB) GetTrips(medallion []string, date string) (map[string]int, error) {
	datStr := date + "%"
	result := make(map[string]int)
	var count int
	for i := 0; i < len(medallion); i++ {
		count = 0
		err := db.QueryRow("select count(*) from cab_trip_data where medallion IN(?) and pickup_datetime LIKE(?)", medallion[i], datStr).Scan(&count)
		if err != nil {
			log.Println("Query failed for medallion: ", medallion[i], err)
			return nil, err
		}
		result[medallion[i]] = count
	}

	return result, nil
}
