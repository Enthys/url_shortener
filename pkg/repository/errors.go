package repository

import "fmt"

type ErrorLinkNotFound struct {
}

func (ErrorLinkNotFound) Error() string {
	return "could not find link"
}

type ErrorIDNotFound struct {
}

func (ErrorIDNotFound) Error() string {
	return "could not find path of link"
}

type ErrorLinkSaveFailure struct {
	Err error
}

func (e ErrorLinkSaveFailure) Error() string {
	return fmt.Sprintf("failed to store the link. Error %s", e.Err.Error())
}

type ErrorIDSaveFailure struct {
	Err error
}

func (e ErrorIDSaveFailure) Error() string {
	return fmt.Sprintf("failed to store the path. Error %s", e.Err.Error())
}
