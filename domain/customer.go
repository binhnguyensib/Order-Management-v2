package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	Id       bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Password string        `json:"-"`
	Phone    string        `json:"phone"`
}

type CustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type CustomerRegiser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type CustomerLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Customer) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashedPassword)
	return nil
}

func (c *Customer) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	return err == nil
}
