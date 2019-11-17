package database

import (
	"database/sql"
	"encoding/json"
)

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password,omitempty"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

func (u *User) ToJson() ([]byte, error){

	jsonString, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return jsonString, nil

}

func (u *User) Save(db *sql.DB) (*User, error){

	stmt := MYSQL_INSERT_USER

	result, err := db.Exec(stmt,u.Email, u.Password, u.Firstname, u.Lastname)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{int(id), u.Email, u.Password, u.Firstname, u.Lastname}, nil

}

func (u *User) GetById(db *sql.DB) (*User, error) {
	stmt := MYSQL_SELECT_USER

	user := &User{}
	row := db.QueryRow(stmt, u.Id)
	err := row.Scan(&user.Id, &user.Email, &user.Firstname, &user.Lastname)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (u *User) GetByCredentials(db *sql.DB) (*User, error) {
	stmt := MYSQL_SELECT_USER_BY_CREDENTIALS

	user := &User{}
	row := db.QueryRow(stmt, u.Email, u.Password)
	err := row.Scan(&user.Id, &user.Email, &user.Firstname, &user.Lastname)
	if err != nil {
		return nil, err
	}
	return user, nil

}