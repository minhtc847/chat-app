package data

import (
	"BE-chat-app/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var AnonymousProfile = &Profile{}

type Profile struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	Gender    bool      `json:"gender"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Activated bool      `json:"activated"`
}
type ProfileModel struct {
	DB *sql.DB
}

func (u *Profile) IsAnonymous() bool {
	return u == AnonymousProfile
}

type password struct {
	plaintext *string //Maximum length of 72 bytes, use pointer to hide password
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateProfile(v *validator.Validator, profile *Profile) {
	v.Check(profile.Name != "", "name", "must be provided")
	v.Check(len(profile.Name) <= 500, "name", "must not be more than 500 bytes long")
	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, profile.Email)
	if profile.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *profile.Password.plaintext)
	}
	if profile.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (m ProfileModel) Insert(profile *Profile) error {
	query := `
		INSERT INTO profile (email, name, image_url, gender, password_hash,activated)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	args := []interface{}{profile.Email, profile.Name, profile.ImageURL, profile.Gender, profile.Password.hash, profile.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&profile.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "profiles_unique_email"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
func (m ProfileModel) GetByEmail(email string) (*Profile, error) {
	query := `
		SELECT id, email, name, image_url, gender, password_hash, created_at, updated_at, activated
		FROM profile
		WHERE email = $1`
	var profile Profile
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&profile.ID, &profile.Email, &profile.Name, &profile.ImageURL, &profile.Gender, &profile.Password.hash, &profile.CreatedAt, &profile.UpdatedAt, &profile.Activated)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &profile, nil
}
func (m ProfileModel) Update(profile *Profile) error {
	query := `
		UPDATE profile
		SET email = $1, name = $2, image_url = $3, gender = $4, password_hash = $5, updated_at = $6, activated = $7
		WHERE id = $8`
	args := []interface{}{profile.Email, profile.Name, profile.ImageURL, profile.Gender, profile.Password.hash, profile.UpdatedAt, profile.Activated, profile.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "profile_unique_email"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
func (m ProfileModel) Get(id uuid.UUID) (*Profile, error) {
	fmt.Println(id.String())
	query := `
        SELECT id, email, name, image_url, gender, password_hash, created_at, updated_at, activated
		FROM profile
		WHERE id = $1`
	var profile Profile
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id.String()).Scan(
		&profile.ID,
		&profile.Email,
		&profile.Name,
		&profile.ImageURL,
		&profile.Gender,
		&profile.Password.hash,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &profile, nil
}
