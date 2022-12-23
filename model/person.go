package model

import "time"

type PersonCreate struct {
	ID         string    `json:"id"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Password   string    `json:"password"`
	Date       time.Time `json:"date"`
}

type Person struct {
	ID         string `json:"id"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type TokenResp struct {
	Data  Person `json:"data"`
	Token Token  `json:"token"`
}

type GetPeople struct {
	GivenName string `json:"givenName"`
	Email     string `json:"email"`
}

type Page struct {
	Data    []GetPeople `json:"data"`
	Matches int64       `json:"matches"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Status bool   `json:"status"`
	Token  string `json:"token"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	ID      string `json:"id"`
}
