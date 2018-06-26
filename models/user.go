package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	ID uuid.UUID
}

type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CreatedAt      time.Time `json:"-" db:"created_at"`
	UpdatedAt      time.Time `json:"-" db:"updated_at"`
	Name           string    `json:"name" db:"name"`
	Email          string    `json:"email" db:"email"`
	PasswordDigest string    `json:"-" db:"password_digest"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u User) HashPwd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u User) Authenticate(password string) bool {
	digest := []byte(u.PasswordDigest)
	err := bcrypt.CompareHashAndPassword(digest, []byte(password))
	return err == nil
}
