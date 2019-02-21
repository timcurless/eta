package user

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	ErrMissingField = "Error: missing %v"
)

type User struct {
	ID        string `json:"" gorm:"primary_key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RewardsID string `json:"rewards_id" gorm:"index"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Salt      string `json:"-"`
}

func New() User {
	u := User{}
	u.NewSalt()
	return u
}

func (u *User) Validate() error {
	if u.FirstName == "" {
		return fmt.Errorf(ErrMissingField, "FirstName")
	}
	if u.LastName == "" {
		return fmt.Errorf(ErrMissingField, "LastName")
	}
	if u.Email == "" {
		return fmt.Errorf(ErrMissingField, "Email")
	}
	if u.Password == "" {
		return fmt.Errorf(ErrMissingField, "Password")
	}
	return nil
}

func (u *User) NewSalt() {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}
