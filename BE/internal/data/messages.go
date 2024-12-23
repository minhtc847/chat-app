package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID              uuid.UUID `json:"id"`
	ChannelId       uuid.UUID `json:"channel_id"`
	ConversationId  uuid.UUID `json:"conversation_id"`
	SenderProfileId uuid.UUID `json:"sender_profile_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Deleted         bool      `json:"deleted"`
	FileUrl         string    `json:"file_url"`
	Type            string    `json:"type"`
}

type MessageModel struct {
	DB *sql.DB
}

func (m MessageModel) Insert(message *Message) error {
	query := `INSERT INTO message (conversation_id, sender_profile_id, content, deleted, file_url, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query,
		message.ConversationId,
		message.SenderProfileId,
		message.Content,
		false,
		message.FileUrl,
		message.Type,
	).Scan(&message.ID)
}
func (m MessageModel) InsertToChannel(message *Message) error {
	query := `INSERT INTO message (channel_id, conversation_id, sender_profile_id, content, deleted, file_url, type) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query,
		message.ChannelId,
		message.ConversationId,
		message.SenderProfileId,
		message.Content,
		false,
		message.FileUrl,
		message.Type,
	).Scan(&message.ID)
}
func (m MessageModel) Get(id uuid.UUID) (*Message, error) {
	query := `SELECT id, channel_id, conversation_id, sender_profile_id, content, created_at, updated_at, deleted, file_url, type FROM message WHERE id = $1`
	var message Message
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&message.ID,
		&message.ChannelId,
		&message.ConversationId,
		&message.SenderProfileId,
		&message.Content,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.Deleted,
		&message.FileUrl,
		&message.Type,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &message, nil
}
func (m MessageModel) Update(message *Message) error {
	query := `UPDATE message SET conversation_id = $1, sender_profile_id = $2, content = $3, updated_at = $4, deleted = $5, file_url = $6, type = $7 
               WHERE id = $8`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query,
		message.ConversationId,
		message.SenderProfileId,
		message.Content,
		time.Now(),
		message.Deleted,
		message.FileUrl,
		message.Type,
		message.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
func (m MessageModel) UpdateToChannel(message *Message) error {
	query := `UPDATE message SET channel_id = $1, sender_profile_id = $2, content = $3, updated_at = $4, deleted = $5, file_url = $6,
                   type = $7 where id = $8`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query,
		message.ChannelId,
		message.SenderProfileId,
		message.Content,
		time.Now(),
		message.Deleted,
		message.FileUrl,
		message.Type,
		message.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
func (m MessageModel) GetMessagesAfterTime(conversationID uuid.UUID, timestamp time.Time) ([]Message, error) {
	query := `
		SELECT id, conversation_id, sender_profile_id, content, created_at, updated_at, deleted, file_url, type
		FROM messages
		WHERE conversation_id = $1 AND created_at > $2 AND deleted = false
		ORDER BY created_at DESC
		LIMIT 20
	`
	rows, err := m.DB.Query(query, conversationID, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.ID,
			&message.ConversationId,
			&message.SenderProfileId,
			&message.Content,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.Deleted,
			&message.FileUrl,
			&message.Type,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
