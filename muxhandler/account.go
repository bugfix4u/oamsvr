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

type AccountMux struct {
	name string
	db   *gorm.DB
}

func NewAccountMux(adb *gorm.DB) *AccountMux {
	am := AccountMux{
		name: "AccountMux",
		db:   adb,
	}

	return &am
}

func (am *AccountMux) GetName() string {
	return am.name
}

func (am *AccountMux) InitRouter(router *mux.Router) {
	if router == nil {
		log.Fatal(errors.New("Fatal: null router"))
	}
	router.HandleFunc("/api/v1/accounts", ValidateCaller(am.getAccounts, am.db)).Methods("GET")
	router.HandleFunc("/api/v1/accounts", ValidateCaller(am.createAccount, am.db)).Methods("POST")
	router.HandleFunc("/api/v1/account/{id}", ValidateCaller(am.getAccount, am.db)).Methods("GET")
	router.HandleFunc("/api/v1/account/{id}", ValidateCaller(am.modifyAccount, am.db)).Methods("PUT")
	router.HandleFunc("/api/v1/account/{id}", ValidateCaller(am.deleteAccount, am.db)).Methods("DELETE")
}

func (am *AccountMux) getAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := dbh.GetAccounts(am.db)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, accounts)
}

func (am *AccountMux) createAccount(w http.ResponseWriter, r *http.Request) {
	var account dbh.Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	if err := account.CreateAccount(am.db); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, account)
}

func (am *AccountMux) getAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Account ID")
		return
	}

	account := dbh.Account{}
	account.ManagedObject.ID = id
	if err := account.GetAccount(am.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "Account not found")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, account)
}

func (am *AccountMux) modifyAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Account ID")
		return
	}

	var account dbh.Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	account.ID = id

	if err := account.ModifyAccount(am.db); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, account)
}

func (am *AccountMux) deleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Account ID")
		return
	}

	account := dbh.Account{}
	account.ManagedObject.ID = id
	if err := account.DeleteAccount(am.db); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
