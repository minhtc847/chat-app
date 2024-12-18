package data

import (
	"github.com/google/uuid"
	"time"
)

type Friendship struct {
	ID           uuid.UUID  `json:"id"`
	Requester_ID uuid.UUID  `json:"requester_id" binding:"required"`
	Receiver_ID  uuid.UUID  `json:"receiver_id" binding:"required"`
	Status       string     `json:"status"`
	Created_at   time.Time  `json:"created_at"`
	Updated_at   *time.Time `json:"updated_at"`
}

func (m ProfileModel) GetAllFriends(userID uuid.UUID) (*[]uuid.UUID, error) {
	query := `
		SELECT CASE
			WHEN receiver_id = $1 THEN requester_id
			WHEN requester_id = $1 THEN receiver_id
		END AS friend_id
		FROM friendship
		WHERE (requester_id = $1 OR receiver_id = $1)
		  AND status = 'Accepted'
	`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendIDs []uuid.UUID
	for rows.Next() {
		var friendID uuid.UUID
		if err := rows.Scan(&friendID); err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, friendID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &friendIDs, nil
}

func (m ProfileModel) SendInvite(requesterID, receiverID uuid.UUID) error {
	query := `INSERT INTO friendship(requester_id, receiver_id, status, created_at, updated_at)
			  VALUES($1, $2, $3, $4, $5)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//Change status = "Pending"
	_, err = stmt.Exec(requesterID, receiverID, "Pending", time.Now(), time.Now())
	if err != nil {
		return err
	}
	return err
}

func (m ProfileModel) GetInvite(requesterID, receiverID uuid.UUID) (*Friendship, error) {
	query := `SELECT * FROM friendship
			  WHERE requester_id = $1 AND receiver_id = $2 AND status = 'Pending'`
	row := m.DB.QueryRow(query, requesterID, receiverID)
	var invite Friendship
	err := row.Scan(&invite.ID, &invite.Requester_ID, &invite.Receiver_ID,
		&invite.Status, &invite.Created_at, &invite.Updated_at)
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (m ProfileModel) ConfirmInvite(friendshipID uuid.UUID, status string) error {
	query := `UPDATE friendship
			  SET status = $1, updated_at = $2
			  WHERE id = $3`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, time.Now(), friendshipID)
	return err
}
