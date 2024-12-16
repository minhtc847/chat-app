package data

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
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

	Friends interface {
		GetAllFriends(userID uuid.UUID) (*[]uuid.UUID, error)
		SendInvite(requesterID, receiverID uuid.UUID) error
		GetInvite(requesterID, receiverID uuid.UUID) (*Friendship, error)
		ConfirmInvite(friendshipID uuid.UUID, status string) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Profiles: ProfileModel{DB: db},
		Friends:  ProfileModel{DB: db},
	}
}
