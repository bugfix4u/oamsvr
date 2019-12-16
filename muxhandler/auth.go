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

package muxhandler

import (
	"encoding/json"
	"errors"
	dbh "oamsvr/dbhandler"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type AuthMux struct {
	name string
	db   *gorm.DB
}

func NewAuthMux(adb *gorm.DB) *AuthMux {
	ah := AuthMux{
		name: "AuthMux",
		db:   adb,
	}

	return &ah
}

func (ah *AuthMux) GetName() string {
	return ah.name
}

func (ah *AuthMux) InitRouter(router *mux.Router) {
	if router == nil {
		log.Fatal(errors.New("Fatal: null router"))
	}
	router.HandleFunc("/api/v1/authenticate", ah.authenticateUserHandler).Methods("POST")
}

func (ah *AuthMux) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var login dbh.Login
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	password := login.Password

	if password == "" {
		RespondWithError(w, http.StatusUnauthorized, "Username/Password invalid")
	}

	user := dbh.User{Username: login.Username, Password: login.Password}

	if err := user.GetUserByUsername(ah.db); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if dbh.CheckPasswordHash(password, user.Password) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"password": user.Password,
		})
		tokenString, error := token.SignedString([]byte("H1r3M3N0W"))
		if error != nil {
			RespondWithError(w, http.StatusUnauthorized, error.Error())
		}
		RespondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
	} else {
		RespondWithError(w, http.StatusUnauthorized, "Username/Password invalid")
	}
}