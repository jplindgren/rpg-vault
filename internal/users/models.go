package users

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	CreatedAt string   `json:"created_at"`
	Password  password `json:"-"`
	Activated bool     `json:"activated"`
	Version   int      `json:"-"`
}

func NewUser(email, name, createdAt string, password_hash []byte, version int, activated bool) *User {
	return &User{
		Email:     email,
		CreatedAt: createdAt,
		Name:      name,
		Password:  password{hash: password_hash},
		Activated: activated,
		Version:   version,
	}
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	//returns a hash in that format: $2b$[cost]$[22-character salt][31-character hash]
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plainTextPassword
	p.hash = hash

	return nil
}

// Declare a new AnonymousUser variable.
var AnonymousUser = &User{}

// Check if a User instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
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

type Token struct {
	Plaintext string    `json:"token"`
	Email     string    `json:"email"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}
