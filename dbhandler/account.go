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
	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
)

type Account struct {
	ManagedObject
	Url 	  string    			`json:"url" gorm:"unique;not null"`
	Username  string    			`json:"username" gorm:"unique;not null"`
	Password  string    			`json:"password" gorm:"not null"`
	Questions []*SecurityQuestion	`json:"questions" gorm:"foreignkey:AccountID;PRELOAD:true"`
}

type SecurityQuestion struct {
	ManagedObject
	AccountID uuid.UUID	`json:"accountId" gorm:"column:account_uuid"`
	Question string 	`json:"question" gorm:"not null"`
	Answer   string 	`json:"answer" gorm:"not null"`
}

// TableName alters the given name
func (account *Account) TableName() string {
	return "account"
}

// BeforeCreate is a GORM callback function to do operation before the create is called.
func (account *Account) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ID", uuid.New())
	if err != nil {
		return err
	}

	password, err := Encrypt(account.Password)
	if err != nil {
		return err
	}

	err = scope.SetColumn("Password", password)
	if err != nil {
		return err
	}

	return nil
}

func (account *Account) BeforeUpdate(scope *gorm.Scope) error {
	password, err := Encrypt(account.Password)
	if err != nil {
		return err
	}

	err = scope.SetColumn("Password", password)
	if err != nil {
		return err
	}
	return nil
}

func (account *Account) AfterFind() error {
	var password string
	password, err := Decrypt(account.Password)
	if err != nil {
		return err
	}

	account.Password = password

	return nil
}

// TableName alters the given name
func (sq *SecurityQuestion) TableName() string {
	return "questions"
}

// BeforeCreate is a GORM callback function to do operation before the create is called.
func (sq *SecurityQuestion) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ID", uuid.New())
	if err != nil {
		return err
	}
	return nil
}

func GetAccounts(db *gorm.DB) ([]*Account, error) {
	var accounts []*Account
	result := db.Set("gorm:auto_preload", true).Find(&accounts)
	err := DbResults(result)

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return accounts, nil
}

func (account *Account) GetAccount(db *gorm.DB) error {
	result := db.Set("gorm:auto_preload", true).Where("uuid = ?", account.ID).Find(account)
	return DbResults(result)
}

func (account *Account) CreateAccount(db *gorm.DB) error {
	result := db.Create(account)
	return DbResults(result)
}

func (account *Account) ModifyAccount(db *gorm.DB) error {
	result := db.Save(account)
	return DbResults(result)
}

func (account *Account) DeleteAccount(db *gorm.DB) error {
	result := db.Delete(account)
	return DbResults(result)
}