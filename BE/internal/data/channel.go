package data

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name" binding:"required"`
	Type       string     `json:"type"`
	Created_by uuid.UUID  `json:"created_by" binding:"required"`
	Created_at time.Time  `json:"created_at"`
	Updated_at *time.Time `json:"updated_at"`
}

func (m ProfileModel) CreateChannel(nameChannel string, userID uuid.UUID) error {
	query := `INSERT INTO channel(name, created_by, created_at, updated_at)
			  VALUES($1, $2, $3, $4)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nameChannel, userID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return err
}

func (m ProfileModel) ExistsChannel(channelID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM channel WHERE id = $1`
	var count int
	err := m.DB.QueryRow(query, channelID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
