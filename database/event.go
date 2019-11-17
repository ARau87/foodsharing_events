package database

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type Event struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	City string `json:"city"`
	Address string `json:"address"`
	GpsLat string `json:"gps_lat,omitempty"`
	GpsLong string `json:"gps_long,omitempty"`
	CreatorId int `json:"creator_id"`
	MaxParticipants int `json:"max_participants"`
	Participants []*User
}

func (e *Event) Save(db *sql.DB) (*Event, error){

	stmt := MYSQL_INSERT_EVENT

	result, err := db.Exec(stmt,e.Name, e.Description, e.City, e.Address, e.GpsLat, e.GpsLong, e.CreatorId, e.MaxParticipants)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Event{int(id), e.Name, e.Description, e.City, e.Address, e.GpsLat, e.GpsLong, e.CreatorId, e.MaxParticipants, e.Participants}, nil

}

func (e *Event) GetById(db *sql.DB) (*Event, error) {

	stmt := MYSQL_SELECT_EVENT_BY_ID
	stmt2 := MYSQL_SELECT_EVENT_PARTICIPANTS_BY_ID

	// Get event data
	event := &Event{}
	row := db.QueryRow(stmt, e.Id)
	err := row.Scan(&event.Id, &event.Name, &event.Description, &event.City, &event.Address, &event.GpsLat, &event.GpsLong, &event.CreatorId, &event.MaxParticipants)
	if err != nil {
		return nil, err
	}

	// Get participants data
	rows, err := db.Query(stmt2, e.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	participants := []*User{}
	for rows.Next() {
		participant := &User{}
		err := rows.Scan(&participant.Id, &participant.Email, &participant.Firstname, &participant.Lastname)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	event.Participants = participants

	return event, nil

}

func (e *Event) ToJson() ([]byte, error){

	jsonString, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return jsonString, nil

}

func (e *Event) AddParticipant(db *sql.DB, user *User) error {

	if len(e.Participants) >= e.MaxParticipants {
		return errors.New("the event's maximum of participants is exceeded")
	}

	e.Participants = append(e.Participants, user)

	stmt := MYSQL_INSERT_EVENT_USER

	_, err := db.Exec(stmt, user.Id, e.Id)
	if err != nil {
		return err
	}

	return nil

}