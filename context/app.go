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

package context

import  (
	"log"
	"net/http"
	mh "oamsvr/muxhandler"
	dbh "oamsvr/dbhandler"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type AppContext struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (asc *AppContext) InitializeDB() {
	var err error
	asc.DB = dbh.GetDbConnection()

	if err != nil {
		log.Fatal(err)
	}
	asc.DB.SingularTable(true)
	asc.DB.Debug().AutoMigrate(
		&dbh.Account{},
		&dbh.SecurityQuestion{},
		&dbh.User{},
	)

	asc.DB.Model(&dbh.SecurityQuestion{}).AddForeignKey("account_uuid", "account(uuid)", "CASCADE", "CASCADE")

	// Setup the default super user
	superUser := dbh.User{
		ManagedObject: dbh.ManagedObject{Name: "Admin User", Description: "Default admin (superuser)"},
		Username:      "admin",
		Password:      "H1r3M3N0W",
		Role:          dbh.SuperUser,
		State:         dbh.UserStateEnabled}

	err = asc.DB.Where("username = ?", superUser.Username).FirstOrCreate(&superUser).Error
	if err != nil {
		log.Fatal(err)
	}

}

func (asc *AppContext) InitializeRouter() {
	asc.Router = mux.NewRouter()

	apiHandlers := [...]mh.MuxHandler{mh.NewAuthMux(asc.DB),
		mh.NewAccountMux(asc.DB),
		mh.NewUserMux(asc.DB)}

	//Setting up handlers
	for _, h := range apiHandlers {
		log.Printf("Setting up mux handlers for %s\n", h.GetName())
		h.InitRouter(asc.Router)
	}
}

func (asc *AppContext) Run(port string) *http.Server {
	svr := &http.Server{Addr: port, Handler: asc.Router}

	go func() {
		log.Fatal(svr.ListenAndServeTLS("./oamsvr.crt", "./oamsvr.key"))
	}()

	return svr
}