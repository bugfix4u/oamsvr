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
	"database/sql"
	"encoding/json"
	"errors"
	dbh "oamsvr/dbhandler"

	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)


type UserMux struct {
	name string
	db   *gorm.DB
}

func NewUserMux(adb *gorm.DB) *UserMux {
	uh := UserMux{
		name: "User",
		db:   adb,
	}

	return &uh
}

func (uh *UserMux) GetName() string {
	return uh.name
}

func (uh *UserMux) InitRouter(router *mux.Router) {
	if router == nil {
		log.Fatal(errors.New("Fatal: null router"))
	}
	router.HandleFunc("/api/v1/users", ValidateUserApiCaller(uh.getUsershandler, uh.db)).Methods("GET")
	router.HandleFunc("/api/v1/users", ValidateUserApiCaller(uh.createUser, uh.db)).Methods("POST")
	router.HandleFunc("/api/v1/user/{id}", ValidateUserApiCaller(uh.getUser, uh.db)).Methods("GET")
	router.HandleFunc("/api/v1/user/{id}", ValidateUserApiCaller(uh.modifyUser, uh.db)).Methods("PUT")
	router.HandleFunc("/api/v1/user/{id}", ValidateUserApiCaller(uh.deleteUser, uh.db)).Methods("DELETE")
}

func (uh *UserMux) getUsershandler(w http.ResponseWriter, r *http.Request) {
	users, err := dbh.GetUsers(uh.db)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, users)
}

func (uh *UserMux) createUser(w http.ResponseWriter, r *http.Request) {
	var user dbh.User
	var err error
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	if err = user.CreateUser(uh.db); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, user)
}

func (uh *UserMux) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user := dbh.User{}
	user.ManagedObject.ID = id
	if err := user.GetUser(uh.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}

func (uh *UserMux) modifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	var user dbh.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	user.ManagedObject.ID = id

	if err := user.ModifyUser(uh.db); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}

func (uh *UserMux) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user := dbh.User{}
	user.ManagedObject.ID = id
	if err := user.DeleteUser(uh.db); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}