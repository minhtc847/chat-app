package data

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	Profiles interface {
		Insert(profile *Profile) error
		GetByEmail(email string) (*Profile, error)
		Update(profile *Profile) error
		Get(id uuid.UUID) (*Profile, error)
	}
	Messages interface {
		Insert(message *Message) error
		InsertToChannel(message *Message) error
		Get(id uuid.UUID) (*Message, error)
		Update(message *Message) error
		GetMessagesAfterTime(conversationID uuid.UUID, timestamp time.Time) ([]Message, error)
	}

	Conversations interface {
		Insert(conversation *Conversations) error
		Get(id uuid.UUID) (*Conversations, error)
		GetByProfiles(profileOneId, profileTwoId uuid.UUID) (*Conversations, error)
		Delete(id uuid.UUID) error
	}
	Friends interface {
		GetAllFriends(userID uuid.UUID) (*[]uuid.UUID, error)
		SendInvite(requesterID, receiverID uuid.UUID) error
		GetInvite(requesterID, receiverID uuid.UUID) (*Friendship, error)
		ConfirmInvite(friendshipID uuid.UUID, status string) error
	}

	Channel interface {
		CreateChannel(nameChannel string, userID uuid.UUID) error
		ExistsChannel(channelID uuid.UUID) (bool, error)
	}

	Member interface {
		AddMembers(channelID uuid.UUID, userIDs []uuid.UUID) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Profiles: ProfileModel{DB: db},
		Messages: MessageModel{DB: db},
		Conversations: ConversationModel{
			DB: db,
		},
		Friends: ProfileModel{DB: db},
		Channel: ProfileModel{DB: db},
		Member:  ProfileModel{DB: db},
	}
}
