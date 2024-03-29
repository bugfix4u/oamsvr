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
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbUser 		= "oamuser"
var dbPassword 	= "H1r3M3N0W"
var dbName 		= "oamdb"
var dbHost 		= "oam-postgres.oam-network"
var dbPort 		= "5432"

// Gorm wants to know the type of database
var dialect 	= "postgres"

func GetDbConnectString() string {
	connectString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port= %s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
	return connectString
}

func GetDbConnection() *gorm.DB {
	conn, err := gorm.Open(dialect, GetDbConnectString())
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func DbResults(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}