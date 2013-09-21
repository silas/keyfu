package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Id       bson.ObjectId "_id,omitempty"
	Email    string
	Password string
	Recover  string
}

func (user *User) SetEmail(email string) (err error) {
	if err = validateEmail(email); err == nil {
		user.Email = email
	}

	return err
}

func (user *User) CheckPassword(password string) (err error) {
	if password == "" || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		err = errors.New("Invalid current password.")
	}

	return err
}

func (user *User) SetPassword(password string) (err error) {
	if err = validatePassword(password); err != nil {
		return err
	}

	if bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), 12); err == nil {
		user.Password = string(bytePassword)
		return nil
	}

	return errors.New("Failed to set password.")
}
