package landb

const landbURL = "https://landb.cern.ch/api/"
const devicesURL = "beta/devices/"
const setsURL = "beta/sets/"

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
	Type     string   `json:"type"` // PERSON or EGROUP or RESERVED
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
