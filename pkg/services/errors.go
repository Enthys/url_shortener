package services

import "fmt"

type ErrorFailedToSaveLink struct {
	Err error
}

func (e ErrorFailedToSaveLink) Error() string {
	return fmt.Sprintf("saving of link failed. Error: %s", e.Err)
}

type ErrorFailedRetrieval struct {
	Err error
}

func (e ErrorFailedRetrieval) Error() string {
	return fmt.Sprintf("failed to retrieve record from repository. Error: %s", e.Err.Error())
}

type ErrorLinkNotFound struct {
}

func (e ErrorLinkNotFound) Error() string {
	return "link does not exist"
}
