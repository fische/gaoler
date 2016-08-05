package errors

import "errors"

var (
	//ErrNotValidRepository returns the following message :
	//"This is not a valid VCS repository"
	ErrNotValidRepository = errors.New("This is not a valid VCS repository")
)
