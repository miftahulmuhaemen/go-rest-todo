package business

import "errors"

var (
	//ErrInvalidSpec Error when data given is not valid on update or insert
	ErrInvalidSpec = errors.New("Given spec is not valid")

	//ErrInvalidExistingUsername  Error when username on given data is already exist on data store
	ErrInvalidExistingUsername = errors.New("Username already exist")

	//ErrInvalidID Error when given ID is not the right type
	ErrInvalidID = errors.New("Not an invalid ID")
)
