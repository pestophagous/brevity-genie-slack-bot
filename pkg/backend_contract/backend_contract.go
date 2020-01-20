package backend_contract

import ()

type Api interface {
	GetUsers() ([]UserData, error)
}

type UserData interface {
	ID() string
	HumanReadableName() string
	AtHandle() string
}
