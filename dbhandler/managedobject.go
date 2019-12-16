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
	"github.com/google/uuid"
)

type ManagedObject struct {
	ID          uuid.UUID `json:"id" gorm:"column:uuid;primary_key"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Tags        string    `json:"tags"`
}