package models

import (
	"BE-chat-app/cmd/api/db"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Image_url  string    `json:"image_url"`
	Gender     bool      `json:"gender" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func GetAllProfile() ([]Profile, error) {
	query := `SELECT * FROM profile`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		err := rows.Scan(&profile.ID, &profile.Email, &profile.Name,
			&profile.Image_url, &profile.Gender, &profile.Password,
			&profile.Created_at, &profile.Updated_at)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func GetProfileByID(profileID uuid.UUID) (*Profile, error) {
	query := `SELECT * FROM profile WHERE id = $1`
	row := db.DB.QueryRow(query, profileID)
	var profile Profile
	err := row.Scan(&profile.ID, &profile.Email, &profile.Name,
		&profile.Image_url, &profile.Gender, &profile.Password,
		&profile.Created_at, &profile.Updated_at)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (p *Profile) CreateProfile() error {
	query := `INSERT INTO profile(email, name, image_url, gender, password, created_at, updated_at)
 			  VALUES($1, $2, $3, $4, $5, $6, $7)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Email, p.Name, p.Image_url, p.Gender, p.Password, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return err
}

func (p Profile) UpdateProfile() error {
	query := `UPDATE profile
			  SET email = $1, name = $2, image_url = $3, gender = $4, password = $5, updated_at = $6
			  WHERE id = $7`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Email, p.Name, p.Image_url, p.Gender, p.Password, time.Now(), p.ID)
	return err
}

func (p Profile) DeleteProfile() error {
	query := `DELETE FROM profile
			  WHERE id = $1`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID)
	return err
}
