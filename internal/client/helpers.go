// Copyright (c) Christopher Barnes <christopher.barnes@cern.ch>
// SPDX-License-Identifier: GPL-3.0-or-later

package landb

const landbURL = "https://landb.cern.ch/api/"

type Location struct {
	Building string `json:"building"`
	Floor    string `json:"floor"`
	Room     string `json:"room"`
}

type OperatingSystem struct {
	Family  string `json:"family"`
	Version string `json:"version"`
}

type Contact struct {
	Type     string   `json:"type"`
	Person   Person   `json:"person"`
	EGroup   EGroup   `json:"egroup"`
	Reserved Reserved `json:"reserved"`
}

type Person struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Department string `json:"department"`
	Group      string `json:"group"`
}

type EGroup struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Reserved struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
