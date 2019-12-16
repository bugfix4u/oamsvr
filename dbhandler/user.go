// Online Account Manager
// Copyright (C) 2019  Denny Chambers

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// package main

package dbhandler

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
)

type User struct {
	ManagedObject
	Username string    `json:"username" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"not null"`
	Role     UserRole  `json:"role" gorm:"default:'admin';not null"` //superuser, admin
	Email    string    `json:"email"`
	State    UserState `json:"status" gorm:"default:'enabled'"` //ENABLED, DISABLED, LOCKED
	Lsl      time.Time `json:"lsl"`                             //Last Successfull Login
	LockTime time.Time
}

type UserRole string

const (
	SuperUser UserRole = "superuser"
	Admin     UserRole = "admin"
)

var UserRoleArray = [...]UserRole{SuperUser, Admin}

type UserState string

const (
	UserStateUnknown  UserState = "unknown"
	UserStateEnabled  UserState = "enabled"
	UserStateDisabled UserState = "disabled"
	UserStateLocked   UserState = "locked"
)

var UserStateArray = [...]UserState{UserStateUnknown,
	UserStateEnabled,
	UserStateDisabled,
	UserStateLocked}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TableName alters the given name
func (user *User) TableName() string {
	return "user"
}

// BeforeCreate is a GORM callback function to do operation before the create is called.
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ID", uuid.New())
	if err != nil {
		return err
	}

	var password string
	password, err = HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = scope.SetColumn("Password", password)
	if err != nil {
		return err
	}

	return nil
}

//User DB functions
func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	err := DbResults(result)

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return users, nil
}

func (user *User) GetUser(db *gorm.DB) error {
	result := db.Where("uuid = ?", user.ID).Find(user)
	return DbResults(result)
}

func (user *User) GetUserByUsername(db *gorm.DB) error {
	result := db.Where("username = ?", user.Username).Find(user)
	return DbResults(result)
}

func (user *User) CreateUser(db *gorm.DB) error {
	result := db.Create(user)
	return DbResults(result)
}

func (user *User) ModifyUser(db *gorm.DB) error {
	result := db.Save(user)
	return DbResults(result)
}

func (user *User) DeleteUser(db *gorm.DB) error {
	result := db.Delete(user)
	return DbResults(result)
}

func ModifyUserState(id uuid.UUID, state UserState, db *gorm.DB) error {
	user := User{}
	user.ID = id
	result := db.Model(&user).Update("state", state)
	return DbResults(result)
}

func ModifyLsl(id uuid.UUID, time time.Time, db *gorm.DB) error {
	user := User{}
	user.ID = id
	result := db.Model(&user).Update("lsl", time)
	return DbResults(result)
}

func ModifyLockTime(id uuid.UUID, time time.Time, db *gorm.DB) error {
	user := User{}
	user.ID = id
	result := db.Model(&user).Update("lockTime", time)
	return DbResults(result)
}