package errors

import "errors"

//VCS
var (
	//ErrNotValidRepository returns the following message :
	//"This is not a valid VCS repository"
	ErrNotValidRepository = errors.New("This is not a valid VCS repository")
	//ErrNotRightRemote returns the following message :
	//"This is not the right remote type"
	ErrNotRightRemote = errors.New("This is not the right remote type")
)
