package main

import (
	"fmt"
	"time"
)

type ValidationRecordMock struct {
	DB       map[string]*ValidationRecord
	mockedID int
}

func (m *ValidationRecordMock) Insert(eventName, personName, personEmail, hash string) error {
	_, ok := m.DB[hash]
	if ok {
		return fmt.Errorf("record with hash %s already exists", hash)
	}

	m.mockedID += 1
	m.DB[hash] = &ValidationRecord{
		ID:          m.mockedID,
		EventName:   eventName,
		PersonName:  personName,
		PersonEmail: personEmail,
		Hash:        hash,
		Created:     time.Now(),
	}
	return nil
}

func (m *ValidationRecordMock) GetByHash(hash string) (*ValidationRecord, error) {
	validator, ok := m.DB[hash]
	if !ok {
		return nil, fmt.Errorf("no record found with hash %s", hash)
	}
	return validator, nil
}
