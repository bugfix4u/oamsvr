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
	"fmt"
	dbh "oamsvr/dbhandler"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func ValidateCaller(next http.HandlerFunc, db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("H1r3M3N0W"), nil
				})
				if error != nil {
					RespondWithError(w, http.StatusUnauthorized, error.Error())
					return
				}
				if token.Valid {
					claims := token.Claims.(jwt.MapClaims)
					var username string = claims["username"].(string)
					log.Printf("Validating token for %s\n", username)
					user := dbh.User{Username: username}
					if err := user.GetUserByUsername(db); err != nil {
						RespondWithError(w, http.StatusUnauthorized, err.Error())
						return
					}
					if user.Password == claims["password"].(string) {
						log.Printf("Token is valid, processing next handler")
						next(w, req)
					} else {
						RespondWithError(w, http.StatusUnauthorized, "Invalid authorization token")

					}
				} else {
					RespondWithError(w, http.StatusUnauthorized, "Invalid authorization token")

				}
			}
		} else {
			RespondWithError(w, http.StatusUnauthorized, "An authorization header is required")
		}
	})
}

//ValidateUserApiCaller validates the caller for User API calls that can change //the DB. Non super user accounts can only change their own user info, where
//super user accounts can change the info for any user.
func ValidateUserApiCaller(next http.HandlerFunc, db *gorm.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var id uuid.UUID
		var err error
		authorizationHeader := req.Header.Get("authorization")
		method := req.Method
		if method == "POST" {
			vars := mux.Vars(req)
			id, err = uuid.Parse(vars["id"])
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, err.Error())
			}
		}
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("H1r3M3N0W"), nil
				})
				if err != nil {
					RespondWithError(w, http.StatusUnauthorized, err.Error())
					return
				}
				if token.Valid {
					claims := token.Claims.(jwt.MapClaims)
					username := claims["username"].(string)
					log.Printf("Validating token for %s\n", username)
					user := dbh.User{Username: username}
					if err := user.GetUserByUsername(db); err != nil {
						RespondWithError(w, http.StatusUnauthorized, err.Error())
						return
					}
					if user.Password == claims["password"].(string) {
						log.Printf("Token password is valid")
						if method == "POST" || method == "DELETE" {
							if user.Role == dbh.SuperUser {
								next(w, req)
							} else {
								RespondWithError(w, http.StatusUnauthorized, "Not authorized for this command")
							}
						} else if method == "PUT" { //User is trying to update existing data
							if user.Role == dbh.SuperUser || id == user.ManagedObject.ID {
								next(w, req)
							} else {
								RespondWithError(w, http.StatusUnauthorized, "Not authorized for this command")
							}
						} else {
							next(w, req)
						}
					} else {
						RespondWithError(w, http.StatusUnauthorized, "Invalid authorization token")
					}
				} else {
					RespondWithError(w, http.StatusUnauthorized, "Invalid authorization token")

				}
			}
		} else {
			RespondWithError(w, http.StatusUnauthorized, "An authorization header is required")
		}
	})
}