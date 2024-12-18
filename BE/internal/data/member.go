package data

import (
	"github.com/google/uuid"
	"time"
)

type Member struct {
	ID         uuid.UUID  `json:"id"`
	ChannelID  uuid.UUID  `json:"channel_id" binding:"required"`
	ProfileID  uuid.UUID  `json:"profile_id" binding:"required"`
	Role       string     `json:"role"`
	Created_at time.Time  `json:"created_at"`
	Updated_at *time.Time `json:"updated_at"`
}

func (m ProfileModel) AddMembers(channelID uuid.UUID, userIDs []uuid.UUID) error {
	query := `INSERT INTO member(channel_id, profile_id, created_at, updated_at)
			  VALUES($1, $2, $3, $4)`
	for _, userID := range userIDs {
		_, err := m.DB.Exec(query, channelID, userID, time.Now(), time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}
