package repository

import "fmt"

type ErrorLinkNotFound struct {
}

func (ErrorLinkNotFound) Error() string {
	return "could not find link"
}

type ErrorPathNotFound struct {
}

func (ErrorPathNotFound) Error() string {
	return "could not find path of link"
}

type ErrorLinkSaveFailure struct {
	Err error
}

func (e ErrorLinkSaveFailure) Error() string {
	return fmt.Sprintf("failed to store the link. Error %s", e.Err.Error())
}

type ErrorPathSaveFailure struct {
	Err error
}

func (e ErrorPathSaveFailure) Error() string {
	return fmt.Sprintf("failed to store the path. Error %s", e.Err.Error())
}
