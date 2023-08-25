package main

import (
	"encoding/json"
	"os"
	"time"
)

type level int

const (
	Owner level = iota
	Admin
	Pro
	Basic
)

func (l level) getAttackStats() (maxTime int, maxConcurrents int, cooldown int) {
	switch l {
	case Owner:
		return 9999, 9999, 0
	case Admin:
		return 999, 10, 10
	case Pro:
		return 60, 2, 30
	case Basic:
		return 30, 1, 60
	}
	return 0, 0, 0
}

type User struct {
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	Expire   time.Time `json:"expire"`
	Level    level     `json:"level"`
}

func (user *User) CanAttack(time int) (bool, string) {
	switch user.Level {
	case Owner:
		break
	case Admin:
		break
	case Pro:
		if time > 60 {
			return false, "Max time is 60"
		}
		break
	case Basic:
		if time > 30 {
			return false, "Max time is 30"
		}
	}
	return true, ""
}

func AuthUser(username string, password string) (bool, *User) {
	users := []User{}
	usersFile, err := os.ReadFile("users.json")
	if err != nil {
		return false, nil
	}
	json.Unmarshal(usersFile, &users)
	for _, user := range users {
		if user.Username == username && user.Password == password {
			if user.Expire.After(time.Now()) {
				return true, &user
			}
		}
	}
	return false, nil
}
