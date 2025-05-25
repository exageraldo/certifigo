package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ValidationRecordModelInterface interface {
	Insert(eventName, personName, personEmail, hash string) error
	GetByHash(hash string) (*ValidationRecord, error)
}

type ValidationRecord struct {
	ID          int
	EventName   string
	PersonName  string
	PersonEmail string
	Hash        string
	Created     time.Time
}

type ValidationRecordModel struct {
	DB *sql.DB
}

func (m *ValidationRecordModel) Insert(eventName, personName, personEmail, hash string) error {
	stmt := `INSERT INTO users (event_name, person_name, person_email, hash, created)
	VALUES(?, ?, ?, ?, datetime('now'))`

	_, err := m.DB.Exec(stmt, eventName, personName, personEmail, hash)
	return err
}

func (m *ValidationRecordModel) GetByHash(hash string) (*ValidationRecord, error) {
	stmt := `SELECT id, event_name, person_name, person_email, hash, created
	FROM users WHERE hash = ?`

	row := m.DB.QueryRow(stmt, hash)

	var v ValidationRecord
	err := row.Scan(&v.ID, &v.EventName, &v.PersonName, &v.PersonEmail, &v.Hash, &v.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no record found with hash %s", hash)
		}
		return nil, err
	}
	return &v, nil
}
